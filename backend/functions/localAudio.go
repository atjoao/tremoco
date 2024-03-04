package functions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"music/server/utils"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

func isAudioFile(filename string) bool {
	audioRegex := regexp.MustCompile(`\.(mp3|wav|flac|aac|opus)$`)
	return audioRegex.MatchString(filename)
}

func ffprobeOutput(path string) (*utils.FFProbeOutputResponse, error){
	var output utils.FFProbeOutputResponse
	var err error
	cmd := "ffprobe"
	args := []string{"-v", "quiet", "-print_format", "json", "-show_format", path}

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(out, &output)
	if err != nil {
		log.Println("Error parsing ffprobe output > ", err)
		return nil, err
	}
	return &output, err
}

func ProcessAudioFiles(db *sql.DB) (bool, error) {
	log.Println("Processing audio files")
	folders, err := os.ReadDir("audio")

	if err != nil {
		return false, err
	}

	for _, e := range folders {
		if e.IsDir() {
			fmt.Println("Folder: ", e.Name())
			const checkAlbumSQL = "SELECT id FROM Album WHERE name = $1"
			
			var albumID int
			err := db.QueryRow(checkAlbumSQL, e.Name()).Scan(&albumID)
			if err != nil {
				if err != sql.ErrNoRows{
					log.Panic("Error checking if album exists > ", err)
				}
				const sql string = "INSERT INTO Album(name) VALUES ($1) RETURNING id"
				err = db.QueryRow(sql, e.Name()).Scan(&albumID)
				if err != nil {
					log.Panic("Error sending inserting album into database > ", err)
				}
				log.Println("Album inserted into database with ID:", albumID)
			} else {
				log.Println("Album already exists in database with ID: ", albumID)
			}
			
			ReadFolder(albumID,filepath.Join("audio", e.Name()), db)
		}
	}

	return true, nil
}

func ReadFolder(albumId int, folder string, db *sql.DB) (bool, error) {
    folders, err := os.ReadDir(folder)
    if err != nil {
        return false, err
    }

    for _, e := range folders {
        if e.IsDir() {
            if _, err := ReadFolder(albumId, filepath.Join(folder, e.Name()), db); err != nil {
                return false, err
            }
        } else {
            if isAudioFile(e.Name()) {
                log.Println("Found audio file: ", e.Name())
                path, _ := filepath.Abs(folder)
                fullPath := filepath.Join(path, e.Name())

                var audioId string
                err := db.QueryRow("SELECT id FROM Music WHERE location = $1", fullPath).Scan(&audioId)
                if err != nil {
					if err != sql.ErrNoRows {
                	    log.Println("Error checking existence of audio file in the database:", err)
                    	continue
					}

					output, err := ffprobeOutput(filepath.ToSlash(fullPath))
					if err != nil {
						log.Println("Error processing audio file > ", err)
						continue
					}

					audioId = "local_" + utils.RandString(24)
					duration, _ := strconv.ParseFloat(output.Format.Duration, 64)
					title := output.Format.Tags.Title
					if title == "" {
						title = e.Name()
					}

					sql := "INSERT INTO Music(id, title, author, duration, genre, location) VALUES ($1, $2, $3, $4, $5, $6)"
					_, err = db.Exec(sql, audioId, title, output.Format.Tags.Artist, math.Round(duration), output.Format.Tags.Genre, fullPath)
					if err != nil {
						log.Println("Error inserting audio file into database > ", err)
						continue
					}
					// associate audio with album
					sql = "INSERT INTO Album_Music(album_id, music_id) VALUES ($1, $2)"
					_, err = db.Exec(sql, albumId, audioId)
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


