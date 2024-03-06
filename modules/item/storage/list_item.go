package storage

import (
	"context"

	"github.com/nthduc/rest-api-basic-todo/common"
	"github.com/nthduc/rest-api-basic-todo/modules/item/model"
)

func (s *sqlStore) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {
	var result []model.TodoItem

	db := s.db.Where("status <> ?", "Deleted") // lọc status có giá trị khác với Deleted là được

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := db.Order("id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit).Find(&result).Error; err != nil {

		return nil, err
	}

	return result, nil
}
