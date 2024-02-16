package model

import (
	"errors"

	"github.com/nthduc/rest-api-basic-todo/common"
)

var (
	ErrTitleIsBlank = errors.New("title cannot be blank")
)

type TodoItem struct {
	common.SQLModel
	Title       string      `json:"title" gorm:"column:title`
	Description string      `json:"description" gorm:"column:description`
	Status      *ItemStatus `json:"status" gorm:"column:status`
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
