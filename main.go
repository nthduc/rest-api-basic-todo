package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gin_item "github.com/nthduc/rest-api-basic-todo/modules/item/transport/gin"
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

func main() {

	dsn := "root:021220@tcp(127.0.0.1:3308)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(db)

	r := gin.Default()
	//r.Use(middleware.Recovery())

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
			items.POST("", gin_item.CreateItem(db))
			items.GET("", gin_item.ListItem(db))
			items.GET("/:id", gin_item.GetItem(db))
			items.PATCH("/:id", gin_item.UpdateItem(db))
			items.DELETE("/:id", gin_item.DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":3000")
}
