package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// `id` int NOT NULL AUTO_INCREMENT,
// `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
// `image` json DEFAULT NULL,
// `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
// `status` enum('Doing','Done','Deleted') DEFAULT 'Doing',
// `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
// `created_at` datetime DEFAULT CURRENT_TIMESTAMP,

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func main() {
	fmt.Println("Hello")
	now := time.Now().UTC()
	item := TodoItem{
		Id:          1,
		Title:       "This is Item 1",
		Description: "This is Item 1",
		Status:      "Doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	jsonData, err := json.Marshal(item)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

	jsonStr := `{"id":1,"title":"This is Item 1","description":"This is Item 1","status":"Doing","created_at":"2024-02-15T10:06:38.1105894Z","updated_at":null}`

	var item2 TodoItem

	if err := json.Unmarshal([]byte(jsonStr), &item2); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item2)
}
