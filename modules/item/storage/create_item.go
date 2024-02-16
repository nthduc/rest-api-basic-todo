package storage

import (
	"context"

	"github.com/nthduc/rest-api-basic-todo/modules/item/model"
)

func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}
