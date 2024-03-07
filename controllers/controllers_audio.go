package controllers

import (
	"log"
	"music/server/functions"
	structs "music/server/utils"

	"github.com/gin-gonic/gin"
)

// TODO: add local audio support for streams
func VideoDataStream(ctx *gin.Context) {
	videoId := ctx.Query("id")
	includeVideo := ctx.Query("videos")
	complete := ctx.Query("complete")

	var includeVideoBool bool
	if includeVideo == "true"{
		includeVideoBool = true
	} else {
		includeVideoBool = false
	}

	if len(videoId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "id" query param or empty`,
		})
		return
	}
	response, metas, err := functions.VideoMeta(videoId, includeVideoBool)
	
	if complete == "true" {
		ctx.PureJSON(200, gin.H{
			"status": "OK",
			"data": response,
		})
		return
	}
	if err != nil {
		log.Printf("Error searching for videos: %v", err)
		ctx.JSON(500, gin.H{
			"status": "SERVER_ERROR",
			"message":  "Internal Server Error",
		})
		return
	}

	if len(metas) == 0 {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
		return
	}

	ctx.PureJSON(200, gin.H{
		"status": "OK",
		"data": &structs.VideoMeta{
			VideoId: response.VideoDetails.VideoId,
			Title: response.VideoDetails.Title,
			Author: response.VideoDetails.Author,
			Thumbnails: response.VideoDetails.Thumbnail.Thumbnails,
			Duration: response.VideoDetails.LengthSeconds,
			Streams: metas,
		},
	})
}

func SearchVideos(ctx *gin.Context) {
	query := ctx.Query("q")
	if len(query) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "q" query param or empty`,
		})
		return
	}

	videos, err := functions.SearchVideo(query)
	if err != nil {
		log.Printf("Error searching for videos: %v", err)
		ctx.JSON(500, gin.H{
			"status": "SERVER_ERROR",
			"message":  "Internal Server Error",
		})
		return
	}

	if len(videos) == 0 {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "OK",
		"videos": videos,
	})

}
