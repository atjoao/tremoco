package controllers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
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
	defer db.Close()

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
	defer db.Close()

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
	// Remove defer db.Close() as we're using a connection pool

	userId := sessions.Default(ctx).Get("userId").(int)
	playlistId := ctx.PostForm("playlistId")
	audioId := ctx.PostForm("audioId")

	if playlistId == "" || audioId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId or audioId is empty",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Failed to start transaction: " + err.Error(),
		})
		return
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM playlists WHERE userId = ? AND id = ?)", userId, playlistId).Scan(&exists)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking playlist ownership: " + err.Error(),
		})
		return
	}
	if !exists {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM playlist_music WHERE playlist_id = ? AND music_id = ?)", playlistId, audioId).Scan(&exists)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking audio in playlist: " + err.Error(),
		})
		return
	}

	var result sql.Result
	if exists {
		result, err = tx.Exec("DELETE FROM playlist_music WHERE playlist_id = ? AND music_id = ?", playlistId, audioId)
	} else {
		result, err = tx.Exec("INSERT INTO playlist_music (playlist_id, music_id) VALUES (?, ?)", playlistId, audioId)
	}

	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error updating playlist: " + err.Error(),
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking affected rows: " + err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(400, gin.H{
			"status":  "NO_CHANGE",
			"message": "No changes were made to the playlist",
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error committing transaction: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Playlist updated successfully",
	})
}

func GetPlaylist(ctx *gin.Context) {
	db := utils.StartConn()
	defer db.Close()
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
	userId := sessions.Default(ctx).Get("userId").(int)
	queryPlaylistId := ctx.Param("playlistId")

	if queryPlaylistId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId is empty",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("[deleteplaylist/err/begin] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM playlists WHERE userId = ? AND id = ?)", userId, queryPlaylistId).Scan(&exists)
	if err != nil {
		log.Printf("[deleteplaylist/err/check] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking playlist ownership",
		})
		return
	}
	if !exists {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	_, err = tx.Exec("DELETE FROM playlist_music WHERE playlist_id = ?", queryPlaylistId)
	if err != nil {
		log.Printf("[deleteplaylist/err/delete_music] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error deleting playlist music entries",
		})
		return
	}

	_, err = tx.Exec("DELETE FROM playlists WHERE id = ?", queryPlaylistId)
	if err != nil {
		log.Printf("[deleteplaylist/err/delete_playlist] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error deleting playlist",
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[deleteplaylist/err/commit] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error committing transaction",
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
	defer db.Close()
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

	userId := sessions.Default(ctx).Get("userId").(int)
	queryPlaylistId := ctx.Param("playlistId")
	playlistName := strings.TrimSpace(ctx.PostForm("playlistName"))
	playlistImage := strings.TrimSpace(ctx.PostForm("playlistImage"))

	if queryPlaylistId == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "playlistId is empty",
		})
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("[editplaylist/err] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback() // Will be a no-op if tx.Commit() is called

	// Check playlist ownership within the transaction
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM playlists WHERE userId = ? AND id = ?)", userId, queryPlaylistId).Scan(&exists)
	if err != nil {
		log.Printf("[editplaylist/err] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking playlist ownership",
		})
		return
	}
	if !exists {
		ctx.JSON(403, gin.H{
			"status":  "FORBIDDEN",
			"message": "Playlist does not belong to user",
		})
		return
	}

	var updateSQL string
	var args []interface{}

	if playlistName == "" {
		updateSQL = "UPDATE playlists SET image = ? WHERE id = ?"
		args = []interface{}{sql.NullString{String: playlistImage, Valid: playlistImage != ""}, queryPlaylistId}
	} else {
		updateSQL = "UPDATE playlists SET name = ?, image = ? WHERE id = ?"
		args = []interface{}{playlistName, sql.NullString{String: playlistImage, Valid: playlistImage != ""}, queryPlaylistId}
	}

	result, err := tx.Exec(updateSQL, args...)
	if err != nil {
		log.Printf("[editplaylist/err] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error updating playlist",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[editplaylist/err] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error checking affected rows",
		})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(400, gin.H{
			"status":  "NO_CHANGE",
			"message": "No changes were made to the playlist",
		})
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[editplaylist/err] %v", err)
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error committing transaction",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "Playlist updated successfully",
	})
}
