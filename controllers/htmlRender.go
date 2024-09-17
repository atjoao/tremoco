package controllers

import (
	"encoding/base64"
	"log"
	"tremoco/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Sidebar(ctx *gin.Context) {
	db := utils.StartConn()
	defer db.Close()

	var userId int = sessions.Default(ctx).Get("userId").(int)

	const sql string = "SELECT id, name, image FROM playlists WHERE userId = $1"

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

		err = rows.Scan(&playlist.PlaylistId, &playlist.PlaylistName, &playlist.PlaylistImage)
		if err != nil {
			log.Println("Error on playlist scan > ", err, "for user id > ", userId)
			continue
		}

		if IsDomainAllowed(playlist.PlaylistImage.String) {
			playlist.PlaylistImage.String = "/api/proxy?url=" + base64.StdEncoding.EncodeToString([]byte(playlist.PlaylistImage.String))
		}

		if playlist.PlaylistImage.String == "" {
			playlist.PlaylistImage.String = "/assets/images/default_album.png"
		}

		playlists = append(playlists, playlist)
	}

	ctx.HTML(200, "sidebar.tmpl", gin.H{
		"playlists": playlists,
	})
}
