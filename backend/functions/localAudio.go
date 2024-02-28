package functions

import (
	"database/sql"
	"log"
	"os"
	"regexp"
)

func isAudioFile(filename string) bool {
	audioRegex := regexp.MustCompile(`\.(mp3|wav|flac|aac|opus)$`)
	return audioRegex.MatchString(filename)
}

func ProcessAudioFiles(db *sql.DB) (bool, error) {
	log.Println("Processing audio files");
	folders, err := os.ReadDir("audio")

	if err != nil {
		return false, err
	}

	for _, e := range folders {
		if e.IsDir() {
			ReadFolder("audio/" + e.Name(), db)
		}
	}

	return true, nil
}

func ReadFolder(folder string, db *sql.DB) (bool, error) {
	folders, err := os.ReadDir(folder)

	if err != nil {
		return false, err
	}

	for _, e := range folders {
		if e.IsDir() {
			ReadFolder(folder + "/" + e.Name(), db)
		} else {
			if isAudioFile(e.Name()) {
				log.Println("Found audio file: ", e.Name())
				//audioId := "local_"+utils.RandString(24)
				//const sql string = "INSERT INTO Music(id, title, author, duration) VALUES ($1, $2, $3, $4)"
				//db.Exec(sql, audioId, )
			}
		}
	}
	return true, nil
}

