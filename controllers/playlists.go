package controllers

import (
	"database/sql"
	"fmt"
	"music/functions"
	"music/utils"
	"os"
	"regexp"
	"strings"

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

func ChangePlaylist(ctx *gin.Context) {
	var db *sql.DB = utils.StartConn()
	var userId int = sessions.Default(ctx).Get("userId").(int)

	var playlistId string = ctx.PostForm("playlistId")
	var audioId string = ctx.PostForm("audioId")

	if playlistId == "" || audioId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId or audioId is empty",
		})
		return
	}

	// check if playlist is from user
	var sql string = "SELECT id FROM playlists WHERE userId = $1 AND id = $2"
	rows, err := db.Query(sql, userId, playlistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	sql = "SELECT music_id FROM playlist_music WHERE playlist_id = $1 AND music_id = $2"
	chkrows, err := db.Query(sql, playlistId, audioId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	if chkrows.Next() {
		sql = "DELETE FROM playlist_music WHERE playlist_id = $1 AND music_id = $2"
		_, err = db.Exec(sql, playlistId, audioId)
		if err != nil {
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while removing music from playlist",
			})
			return
		}
	} else {
		sql = "INSERT INTO playlist_music (playlist_id, music_id) VALUES ($1, $2)"
		_, err = db.Exec(sql, playlistId, audioId)
		if err != nil {
			fmt.Println("[erro! changeplaylist] ", err)
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while adding music to playlist",
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Music added to playlist",
	})

}

func GetPlaylist(ctx *gin.Context) {
	var db *sql.DB = utils.StartConn()
	var userId int = sessions.Default(ctx).Get("userId").(int)
	var queryPlaylistId string = ctx.Param("playlistId")
	idRegex := regexp.MustCompile(`local-[^\s]+`)

	if queryPlaylistId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId is empty",
		})
		return
	}

	var sql string = "SELECT id, name, image FROM playlists WHERE userId = $1 AND id = $2"
	rows, err := db.Query(sql, userId, queryPlaylistId)
	if err != nil {
		fmt.Println("[GetPlaylist/err] ", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	var playlistName string
	var playlistId int
	var playlistImage string
	rows.Scan(&playlistId, &playlistName, &playlistImage)

	sql = "SELECT music_id FROM playlist_music WHERE playlist_id = $1"
	playlistRows, err := db.Query(sql, playlistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	var playlistMusics []utils.VideoMeta

	for playlistRows.Next() {
		var musicId string
		err := playlistRows.Scan(&musicId)
		if err != nil {
			continue
		}

		if idRegex.MatchString(musicId) {
			response := functions.LocalVideoMeta(musicId)
			if response == nil {
				continue
			}

			playlistMusics = append(playlistMusics, *response)
		} else if os.Getenv("INCLUDE_YOUTUBE") == "true" {

			response, metas, err := functions.VideoMeta(musicId, false)

			if len(metas) == 0 {
				continue
			}

			if err != nil {
				fmt.Println("[GETPLAYLIST/err] ", err)
			}

			author := strings.Split(response.VideoDetails.Author, "-")

			music := &utils.VideoMeta{
				VideoId:    response.VideoDetails.VideoId,
				Title:      response.VideoDetails.Title,
				Author:     strings.Trim(author[0], " "),
				Thumbnails: response.VideoDetails.Thumbnail.Thumbnails,
				Duration:   response.VideoDetails.LengthSeconds,
				Streams:    metas,
			}

			playlistMusics = append(playlistMusics, *music)
		}
	}

	ctx.JSON(200, gin.H{
		"status": "OK",
		"playlist": &utils.Playlist{
			PlaylistImage: playlistImage,
			PlaylistId:    playlistId,
			PlaylistName:  playlistName,
			MusicList:     playlistMusics,
		},
	})

}

func DeletePlaylist(ctx *gin.Context) {
	var db *sql.DB = utils.StartConn()

	var userId int = sessions.Default(ctx).Get("userId").(int)
	var queryPlaylistId string = ctx.Param("playlistId")

	if queryPlaylistId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId is empty",
		})
		return
	}

	var sql string = "SELECT id FROM playlists WHERE userId = $1 AND id = $2"
	rows, err := db.Query(sql, userId, queryPlaylistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	sql = "DELETE FROM playlist_music WHERE playlist_id = $1"
	_, err = db.Exec(sql, queryPlaylistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	sql = "DELETE FROM playlists WHERE id = $1"
	_, err = db.Exec(sql, queryPlaylistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Playlist deleted",
	})

}
