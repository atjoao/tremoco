package controllers

import (
	"log"
	"music/functions"
	"music/utils"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

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
	if includeVideo == "true" {
		includeVideoBool = true
	} else {
		includeVideoBool = false
	}

	idRegex := regexp.MustCompile(`local-[^\s]+`)
	if idRegex.MatchString(videoId) {
		response := functions.LocalVideoMeta(videoId)
		if response == nil {
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		}

		ctx.PureJSON(200, gin.H{
			"status": "OK",
			"data": &utils.VideoMeta{
				VideoId:    response.VideoId,
				Title:      response.Title,
				Author:     response.Author,
				Thumbnails: response.Thumbnails,
				Duration:   response.Duration,
				Streams:    response.Streams,
			},
		})

	} else if os.Getenv("INCLUDE_YOUTUBE") == "true" {
		response, metas, err := functions.VideoMeta(videoId, includeVideoBool)

		if complete == "true" {
			ctx.PureJSON(200, gin.H{
				"status": "OK",
				"data":   response,
			})
			return
		}
		if err != nil {
			log.Printf("Error searching for videos: %v", err)
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "Internal Server Error",
			})
			return
		}

		if len(metas) == 0 {
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		}

		author := strings.Split(response.VideoDetails.Author, "-")

		ctx.PureJSON(200, gin.H{
			"status": "OK",
			"data": &utils.VideoMeta{
				VideoId:    response.VideoDetails.VideoId,
				Title:      response.VideoDetails.Title,
				Author:     strings.Trim(author[0], " "),
				Thumbnails: response.VideoDetails.Thumbnail.Thumbnails,
				Duration:   response.VideoDetails.LengthSeconds,
				Streams:    metas,
			},
		})
	} else {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
		return
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
			"status":  "SERVER_ERROR",
			"message": "Internal Server Error",
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
