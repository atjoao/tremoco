package routers

import (
	"log"
	"music/server/functions"

	"github.com/gin-gonic/gin"
)

type video struct {
	title string
	videoId string
	thumbnailUrl string
}

func TestGroup(rg *gin.RouterGroup) {
	rg.GET("/a", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})


	rg.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id");
		ctx.JSON(200, gin.H{
			"message": id,
		})
	})

	// TODO: build cache thing
	rg.GET("/search", func(ctx *gin.Context) {
		query := ctx.Query("q")
		if len(query) <= 0 {
			ctx.JSON(400, gin.H{
				"status": "MISSED_PARAMS",
				"message": "no \"q\" query param or empty",
			})
			return;
		}

		videos,err := functions.SearchVideo(query)
		if err != nil{
			log.Printf("Error searching for videos: %v", err)
			ctx.JSON(500, gin.H{
				"status": "SERVER_ERROR",
				"error": "Internal Server Error",
			})
			return;
		}

		ctx.JSON(200, gin.H{
			"status": "OK",
			"videos": videos,
		})
	})

	// entendi como usar isto ig
	rg.StaticFile("/audio", "./audio/audio.opus")
}