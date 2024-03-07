package controllers

import (
	"log"
	"music/server/functions"
	structs "music/server/utils"
	"regexp"

	"github.com/gin-gonic/gin"
)

func StreamAudio(ctx *gin.Context) {
	audioId := ctx.Param("audioId")
	if len(audioId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "audioId" param or empty`,
		})
		return
	}

	idRegex := regexp.MustCompile(`local_[^\s]+`)
	if idRegex.MatchString(audioId){
		_, location := functions.LocalVideoMeta(audioId)
		if location == ""{
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		} else {
			ctx.File(location)
		}
	} else {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
	}
}

func VideoDataStream(ctx *gin.Context) {
	videoId := ctx.Query("id")
	if len(videoId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "id" query param or empty`,
		})
		return
	}

	includeVideo := ctx.Query("videos")
	complete := ctx.Query("complete")

	var includeVideoBool bool
	if includeVideo == "true"{
		includeVideoBool = true
	} else {
		includeVideoBool = false
	}

	idRegex := regexp.MustCompile(`local_[^\s]+`)
	if idRegex.MatchString(videoId){
		response, _ := functions.LocalVideoMeta(videoId)
		if response == nil{
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		}

		ctx.PureJSON(200, gin.H{
			"status": "OK",
			"data": &structs.VideoMeta{
				VideoId: response.VideoId,
				Title: response.Title,
				Author: response.Author,
				Thumbnails: response.Thumbnails,
				Duration: response.Duration,
				Streams: response.Streams,
			},
		})

	} else {
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
