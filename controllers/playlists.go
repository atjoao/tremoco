package controllers

import (
	"database/sql"
	"fmt"
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

// GetPlaylistsMusic returns all playlists that contains the music_id
// used in GET /api/playlist/add
// post will be used to add music to a playlist

func GetPlaylistsMusic(ctx *gin.Context) {
	var db *sql.DB = utils.StartConn()
	var playlists []utils.PlayList_Music

	var userId int = sessions.Default(ctx).Get("userId").(int)

	audioId := ctx.Param("audioId")
	if len(audioId) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "audioId" param or empty`,
		})
		return
	}

	var sql string = "SELECT id, name FROM playlists WHERE userId = $1"
	rows, err := db.Query(sql, userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
	}

	defer rows.Close()

	for rows.Next() {
		var playlist utils.Playlist
		err := rows.Scan(&playlist.PlaylistId, &playlist.PlaylistName)
		if err != nil {
			continue
		}

		sql = "SELECT music_id FROM playlist_music WHERE playlist_id = $1 AND music_id = $2"
		musicRows, err := db.Query(sql, playlist.PlaylistId, audioId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer musicRows.Close()

		playlists = append(playlists, utils.PlayList_Music{
			PlaylistId:   playlist.PlaylistId,
			PlaylistName: playlist.PlaylistName,
			Exists:       musicRows.Next(),
		})
	}

	ctx.JSON(200, gin.H{
		"status":    "OK",
		"playlists": playlists,
	})
}
