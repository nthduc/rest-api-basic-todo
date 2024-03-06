package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nthduc/rest-api-basic-todo/common"
)

func Recovery() func(*gin.Context) {
	return func(c *gin.Context) {
		log.Println("Recovery")

		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrInternal(err))
				}
				log.Println("Recovery", r)
			}
		}()
		c.Next()
	}
}
