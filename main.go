package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nthduc/rest-api-basic-todo/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// `id` int NOT NULL AUTO_INCREMENT,
// `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
// `image` json DEFAULT NULL,
// `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
// `status` enum('Doing','Done','Deleted') DEFAULT 'Doing',
// `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
// `created_at` datetime DEFAULT CURRENT_TIMESTAMP,

type ItemStatus int

const (
	ItemStatusDoing ItemStatus = iota
	ItemStatusDone
	ItemStatusDeleted
)

var allItemStatuses = [3]string{"Doing", "Done", "Deleted"}

func (item *ItemStatus) String() string {
	return allItemStatuses[*item]
}

func parseStr2ItemStatus(s string) (ItemStatus, error) {
	for i := range allItemStatuses {
		if allItemStatuses[i] == s {
			return ItemStatus(i), nil
		}
	}

	return ItemStatus(0), errors.New("Invalid Status String")
}

func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return errors.New(fmt.Sprintf("Fail to scan data from sql: %s", value))
	}

	v, err := parseStr2ItemStatus(string(bytes))

	if err != nil {
		return errors.New(fmt.Sprintf("fail to scan data from sql: %s", value))
	}

	*item = v

	return nil
}

func (item *ItemStatus) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return item.String(), nil
}

// ItemStatus --> JSON value
func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	if item == nil {
		return nil, nil
	}
	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

// JSON value ---> ItemStatus
func (item *ItemStatus) UnmarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "") // "Doing"

	itemValue, err := parseStr2ItemStatus(str)

	if err != nil {
		return err
	}

	*item = itemValue

	return nil
}

type TodoItem struct {
	Id          int         `json:"id" gorm:"column:id`
	Title       string      `json:"title" gorm:"column:title`
	Description string      `json:"description" gorm:"column:description`
	Status      *ItemStatus `json:"status" gorm:"column:status`
	CreatedAt   *time.Time  `json:"created_at" gorm:"column:created_at`
	UpdatedAt   *time.Time  `json:"updated_at,omitempty" gorm:"column:updated_at`
}

type TodoItemCreation struct {
	Id          int         `json:"-" gorm:"column:id`
	Title       string      `json:"title" gorm:"column:title"`
	Description string      `json:"description" gorm:"column:description"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title"`
	Description *string `json:"description" gorm:"column:description"`
	Status      *string `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }

func main() {

	dsn := "root:021220@tcp(127.0.0.1:3308)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(db)

	// now := time.Now().UTC()
	// item := TodoItem{
	// 	Id:          1,
	// 	Title:       "This is Item 1",
	// 	Description: "This is Item 1",
	// 	Status:      "Doing",
	// 	CreatedAt:   &now,
	// 	UpdatedAt:   nil,
	// }

	// jsonData, err := json.Marshal(item)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(jsonData))

	// jsonStr := `{"id":1,"title":"This is Item 1","description":"This is Item 1","status":"Doing","created_at":"2024-02-15T10:06:38.1105894Z","updated_at":null}`

	// var item2 TodoItem

	// if err := json.Unmarshal([]byte(jsonStr), &item2); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(item2)

	r := gin.Default()

	//CRUD: Create, Read, Update, Delete
	//POST: /v1/items (Create a new item)
	//GET: /v1/items (List item) /v1/items?page=1 || /v1/items?cursor=nthduc
	//GET: /v1/items/:id (get item detail by id)
	//(PUT || PATCH): /v1/items/:id (update a item by id)
	//DELETE /v1/items/:id (delete item detail by id)

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":3000")
}

func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemCreation
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))

	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItem

		// /v1/items/1
		id, err := strconv.Atoi(c.Param("id")) // "id"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// data.Id = id
		// db.First(&data)
		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemUpdate

		// /v1/items/1
		id, err := strconv.Atoi(c.Param("id")) // "id"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// data.Id = id
		// db.First(&data)
		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		// /v1/items/1
		id, err := strconv.Atoi(c.Param("id")) // "id"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// data.Id = id
		// db.First(&data)
		if err := db.Where("id = ?", id).Updates(map[string]interface{}{
			"status": "Deleted",
		}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging
		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Process()

		var result []TodoItem

		db = db.Where("status <> ?", "Deleted") // lọc status có giá trị khác với Deleted là được

		if err := db.Table(TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		if err := db.Order("id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit).Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})

	}
}
