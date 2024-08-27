package controllers

import (
	"database/sql"
	"log"
	"music/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Sidebar(ctx *gin.Context) {
	var db *sql.DB = utils.StartConn()
	var userId int = sessions.Default(ctx).Get("userId").(int)

	const sql string = "SELECT id, name FROM playlists WHERE userId = $1"

	rows, err := db.Query(sql, userId)
	if err != nil {
		log.Println("Error on playlists query > ", err)
		ctx.HTML(500, "sidebar.tmpl", gin.H{
			"playlists": nil,
		})
		return
	}

	defer rows.Close()

	var playlists []utils.Playlist

	for rows.Next() {
		var playlist utils.Playlist
		err = rows.Scan(&playlist.PlaylistId, &playlist.PlaylistName)
		if err != nil {
			log.Println("Error on playlist scan > ", err, "for user id > ", userId)
			continue
		}

		playlists = append(playlists, playlist)
	}

	ctx.HTML(200, "sidebar.tmpl", gin.H{
		"playlists": playlists,
	})
}
