package controllers

import (
	"music/functions"
	"regexp"

	"github.com/gin-gonic/gin"
)

func GetAudioCover(ctx *gin.Context) {
	audioId := ctx.Param("audioId")
	if len(audioId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "audioId" param or empty`,
		})
		return
	}
	idRegex := regexp.MustCompile(`local-[^\s]+`)
	if idRegex.MatchString(audioId) {
		music := functions.LocalVideoMeta(audioId)
		if music == nil {
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		} else {
			ctx.File(music.Cover)
		}
	} else {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
	}
}

func StreamAudio(ctx *gin.Context) {
	audioId := ctx.Param("audioId")
	if len(audioId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "audioId" param or empty`,
		})
		return
	}

	idRegex := regexp.MustCompile(`local-[^\s]+`)
	if idRegex.MatchString(audioId) {
		music := functions.LocalVideoMeta(audioId)

		if music == nil || music.Location == "" {
			ctx.JSON(404, gin.H{
				"status": "NOTHING_FOUND",
			})
			return
		}

		ctx.File(music.Location)
	} else {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
	}
}
