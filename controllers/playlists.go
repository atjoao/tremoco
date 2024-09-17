package controllers

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tremoco/functions"
	"tremoco/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CreatePlaylist(ctx *gin.Context) {
	var err error
	db := utils.StartConn()

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
	db := utils.StartConn()
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
	db := utils.StartConn()
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

	defer rows.Close()

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

	defer chkrows.Close()

	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while executing this action",
		})
		return
	}

	if chkrows.Next() {
		sql = "DELETE FROM playlist_music WHERE playlist_id = $1 AND music_id = $2"
		_, err = tx.Exec(sql, playlistId, audioId)
		if err != nil {
			tx.Rollback()
			fmt.Println("[erro! changeplaylist] ", err)
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while removing music from playlist",
			})
			return
		}
	} else {
		sql = "INSERT INTO playlist_music (playlist_id, music_id) VALUES ($1, $2)"
		_, err = tx.Exec(sql, playlistId, audioId)
		if err != nil {
			tx.Rollback()
			fmt.Println("[erro! changeplaylist] ", err)
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while adding music to playlist",
			})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while executing this action",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Music added to playlist",
	})

}

func GetPlaylist(ctx *gin.Context) {
	db := utils.StartConn()
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
	defer rows.Close()

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	var playlist utils.Playlist
	rows.Scan(&playlist.PlaylistId, &playlist.PlaylistName, &playlist.PlaylistImage)

	if IsDomainAllowed(playlist.PlaylistImage.String) {
		playlist.PlaylistImage.String = "/api/proxy?url=" + base64.StdEncoding.EncodeToString([]byte(playlist.PlaylistImage.String))
	}

	sql = "SELECT music_id FROM playlist_music WHERE playlist_id = $1"
	playlistRows, err := db.Query(sql, playlist.PlaylistId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	defer playlistRows.Close()

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

			playlist.MusicList = append(playlist.MusicList, *response)
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

			playlist.MusicList = append(playlist.MusicList, *music)
		}
	}

	ctx.JSON(200, gin.H{
		"status":   "OK",
		"playlist": playlist,
	})
}

func DeletePlaylist(ctx *gin.Context) {
	db := utils.StartConn()

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

	defer rows.Close()

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	sql = "DELETE FROM playlist_music WHERE playlist_id = $1"
	_, err = tx.Exec(sql, queryPlaylistId)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	sql = "DELETE FROM playlists WHERE id = $1"
	_, err = tx.Exec(sql, queryPlaylistId)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	err = tx.Commit()
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

func GetUserPlaylists(ctx *gin.Context) {
	db := utils.StartConn()
	var userId int = sessions.Default(ctx).Get("userId").(int)

	const sql string = "SELECT id, name, image FROM playlists WHERE userId = $1"

	rows, err := db.Query(sql, userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while fetching playlists",
		})
		return
	}

	defer rows.Close()

	var playlists []utils.PlaylistShow

	for rows.Next() {
		var playlist utils.PlaylistShow
		err = rows.Scan(&playlist.PlaylistId, &playlist.PlaylistName, &playlist.PlaylistImage)
		if err != nil {
			continue
		}

		if IsDomainAllowed(playlist.PlaylistImage) {
			playlist.PlaylistImage = "/api/proxy?url=" + base64.StdEncoding.EncodeToString([]byte(playlist.PlaylistImage))
		}

		playlist.PlaylistUrl = "/api/playlist/" + strconv.Itoa(playlist.PlaylistId)

		playlists = append(playlists, playlist)
	}

	ctx.JSON(200, gin.H{
		"status":    "OK",
		"playlists": playlists,
	})
}

func EditPlaylist(ctx *gin.Context) {
	db := utils.StartConn()

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

	defer rows.Close()

	if !rows.Next() {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while deleting playlist",
		})
		return
	}

	var playlistName string = ctx.PostForm("playlistName")
	var playlistImage string = ctx.PostForm("playlistImage")

	if strings.Trim(playlistImage, " ") == "" {
		// set to null on mysql
		playlistImage = ""
	}

	if strings.Trim(playlistName, " ") == "" {
		sql = "UPDATE playlists SET image = $1 WHERE id = $2"
		_, err = tx.Exec(sql, playlistImage, queryPlaylistId)
		if err != nil {
			tx.Rollback()
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while updating playlist",
			})
			return
		}
	} else {
		sql = "UPDATE playlists SET name = $1, image = $2 WHERE id = $3"
		_, err = tx.Exec(sql, playlistName, playlistImage, queryPlaylistId)
		if err != nil {
			tx.Rollback()
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "There was an error while updating playlist",
			})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "There was an error while updating playlist",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Playlist updated",
	})

}
