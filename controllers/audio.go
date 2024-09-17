package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"tremoco/functions"
	"tremoco/utils"

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
		} else if music.Cover != "" {
			if _, err := os.Stat(music.Cover); os.IsNotExist(err) {
				ctx.FileFromFS("images/no-cover.png", http.FS(utils.Assets))
				return
			}
			ctx.File(music.Cover)
		} else {
			ctx.FileFromFS("images/no-cover.png", http.FS(utils.Assets))
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

		ctx.Header("Content-Type", music.Streams[0].MimeType)
		ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(music.Location))
		ctx.File(music.Location)
	} else {
		ctx.JSON(404, gin.H{
			"status": "NOTHING_FOUND",
		})
	}
}
