package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestGroup(rg *gin.RouterGroup) {
	rg.GET("/a", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "hello world")
	})
}