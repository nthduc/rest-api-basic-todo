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

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemUpdate

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

		store := storage.NewSQLStore(db)
		biz := business.NewUpdateItemBiz(store)

		if err := biz.UpdateItemById(c.Request.Context(), id, &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
