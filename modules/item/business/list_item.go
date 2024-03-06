package business

import (
	"context"

	"github.com/nthduc/rest-api-basic-todo/common"
	"github.com/nthduc/rest-api-basic-todo/modules/item/model"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) // condion
}

type listItemBiz struct {
	store ListItemStorage
}

func NewListItemBiz(store ListItemStorage) *listItemBiz {
	return &listItemBiz{store: store}
}

func (biz *listItemBiz) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {

	data, err := biz.store.ListItem(ctx, filter, paging)

	if err != nil {
		return nil, err
	}

	return data, nil
}
