package gin_item

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nthduc/rest-api-basic-todo/common"
	"github.com/nthduc/rest-api-basic-todo/modules/item/business"
	"github.com/nthduc/rest-api-basic-todo/modules/item/model"
	"github.com/nthduc/rest-api-basic-todo/modules/item/storage"
	"gorm.io/gorm"
)

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItem

		// /v1/items/1
		id, err := strconv.Atoi(c.Param("id")) // "id"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)
		biz := business.NewGetItemBiz(store)

		biz.GetItemById(c.Request.Context(), id)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))

	}
}
