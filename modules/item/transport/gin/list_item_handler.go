package gin_item

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nthduc/rest-api-basic-todo/common"
	"github.com/nthduc/rest-api-basic-todo/modules/item/business"
	"github.com/nthduc/rest-api-basic-todo/modules/item/model"
	"github.com/nthduc/rest-api-basic-todo/modules/item/storage"
	"gorm.io/gorm"
)

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

		var filter model.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)
		biz := business.NewListItemBiz(store)

		result, err := biz.ListItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))

	}
}
