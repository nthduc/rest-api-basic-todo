package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

type TodoItem struct {
	Id          int        `json:"id" gorm:"column:id`
	Title       string     `json:"title" gorm:"column:title`
	Description string     `json:"description" gorm:"column:description`
	Status      string     `json:"status" gorm:"column:status`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at`
}

type TodoItemCreation struct {
	Id          int    `json:"-" gorm:"column:id`
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	//Status      string `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

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
			items.GET("")
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id")
			items.DELETE("/:id")
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

		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})

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

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
