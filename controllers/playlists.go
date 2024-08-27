package controllers

import (
	"database/sql"
	"music/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreatePlaylist(ctx *gin.Context) {
	var err error
	var db *sql.DB = utils.StartConn()

	var playlistName string = ctx.PostForm("name")

	var userId int = sessions.Default(ctx).Get("userId").(int)
	if playlistName == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "Playlist name is empty",
		})
		return
	}

	const sql string = "INSERT INTO playlists (name, userId) VALUES ($1, $2) RETURNING id"
	var playlistId int

	err = db.QueryRow(sql, playlistName, userId).Scan(&playlistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while creating this playlist",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Playlist created",
	})
}
