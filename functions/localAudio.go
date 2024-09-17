package functions

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"tremoco/utils"
)

func isAudioFile(filename string) bool {
	audioRegex := regexp.MustCompile(`\.(mp3|wav|flac|aac|opus)$`)
	return audioRegex.MatchString(filename)
}

func findCover(dirPath string) (string, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			switch strings.ToLower(entry.Name()) {
			case "cover.jpg", "cover.png", "cover.jpeg":
				log.Println("Cover found: ", entry.Name())
				return filepath.Join(dirPath, entry.Name()), nil
			}
		}
	}
	return "", nil
}

func FfprobeOutput(path string) (*utils.FFProbeOutputResponse, error) {
	var output utils.FFProbeOutputResponse
	var err error
	cmd := "ffprobe"
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", path}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(out, &output)
	if err != nil {
		log.Println("Error parsing ffprobe output > ", err)
		return nil, err
	}
	return &output, err
}

func ProcessAudioFiles() (bool, error) {
	db := utils.StartConn()
	defer db.Close()

	var err error

	log.Println("Processing audio files")
	folders, err := os.ReadDir("audio")

	if err != nil {
		return false, err
	}

	for _, e := range folders {
		if e.IsDir() {
			log.Println("Folder: ", e.Name())
			currPath := filepath.Join("audio", e.Name())
			cover, err := findCover(currPath)

			if err != nil {
				log.Println("Error finding cover > ", err)
			}

			if cover == "" {
				log.Println("No cover found for album: ", e.Name())
			}

			tx, err := db.Begin()
			if err != nil {
				return false, err
			}

			const checkAlbumSQL = "SELECT id FROM Album WHERE name = $1"

			var albumID int
			err = tx.QueryRow(checkAlbumSQL, e.Name()).Scan(&albumID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Panic("Error checking if album exists > ", err)
				}
				const sql string = "INSERT INTO Album(name,cover) VALUES ($1,$2) RETURNING id"
				err = tx.QueryRow(sql, e.Name(), cover).Scan(&albumID)
				if err != nil {
					tx.Rollback()
					log.Panic("Error sending inserting album into database > ", err)
				}
				log.Println("Album inserted into database with ID:", albumID)
			} else {
				log.Println("Album already exists in database with ID: ", albumID)
			}

			if success, err := ReadFolder(albumID, currPath, tx); err != nil || !success {
				defer tx.Rollback()
				log.Panic("Error processing folder > ", err)
			} else {
				tx.Commit()
				log.Println("Folder processed successfully")
			}
		}
	}

	return true, nil
}

func ReadFolder(albumId int, folder string, tx *sql.Tx) (bool, error) {

	folders, err := os.ReadDir(folder)
	if err != nil {
		return false, err
	}

	for _, e := range folders {
		if e.IsDir() {
			if _, err := ReadFolder(albumId, filepath.Join(folder, e.Name()), tx); err != nil {
				return false, err
			}
		} else {
			if isAudioFile(e.Name()) {
				log.Println("Found audio file: ", e.Name())
				fullPath := filepath.Join(folder, e.Name())

				var audioId string
				err := tx.QueryRow("SELECT id FROM Music WHERE location = $1", fullPath).Scan(&audioId)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Println("Error checking existence of audio file in the database:", err)
						continue
					}

					output, err := FfprobeOutput(filepath.ToSlash(fullPath))
					if err != nil {
						log.Println("Error processing audio file > ", err)
						continue
					}

					audioId = "local-" + utils.RandString(24)
					duration, _ := strconv.ParseFloat(output.Format.Duration, 64)
					title := output.Format.Tags.Title
					if title == "" {
						title = e.Name()
					}

					sql := "INSERT INTO Music(id, title, author, duration, genre, location) VALUES ($1, $2, $3, $4, $5, $6)"
					_, err = tx.Exec(sql, audioId, title, output.Format.Tags.Artist, math.Round(duration), output.Format.Tags.Genre, fullPath)
					if err != nil {
						log.Println("Error inserting audio file into database > ", err)
						continue
					}
					// associate audio with album
					sql = "INSERT INTO Album_Music(album_id, music_id) VALUES ($1, $2)"
					_, err = tx.Exec(sql, albumId, audioId)
					if err != nil {
						log.Println("Error associating audio file with album > ", err)
						continue
					}
				} else {
					log.Println("Audio file already exists in the database:", e.Name())
					continue
				}

			}
		}
	}
	return true, nil
}

func RemoveMusicFromDb() {
	db := utils.StartConn()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction > ", err)
		return
	}

	locations, err := db.Query("SELECT location FROM Music")
	if err != nil {
		log.Println("Error fetching music locations > ", err)
		return
	}

	defer locations.Close()

	for locations.Next() {
		var location string
		err = locations.Scan(&location)
		if err != nil {
			log.Println("Error scanning location > ", err)
			continue
		}
		if _, err := os.Stat(location); os.IsNotExist(err) {
			log.Println("File not found, removing from database")

			var musicID string
			err = tx.QueryRow("SELECT id FROM Music WHERE location = $1", location).Scan(&musicID)
			if err != nil {
				log.Println("Error getting music id > ", err)
				continue
			}

			_, err = tx.Exec("DELETE FROM Album_Music WHERE music_id = $1", musicID)
			if err != nil {
				log.Println("Error removing music from Album_Music table > ", err)
				continue
			}

			_, err = tx.Exec("DELETE FROM Playlist_Music WHERE music_id = $1", musicID)
			if err != nil {
				log.Println("Error removing music from Playlist_Music table > ", err)
				continue
			}

			_, err = tx.Exec("DELETE FROM Music WHERE id = $1", musicID)
			if err != nil {
				log.Println("Error removing music from Music table > ", err)
				continue
			}

			log.Println("Deleted music with id", musicID)
		}
	}

	tx.Commit()
}
