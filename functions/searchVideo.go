package functions

import (
	"fmt"
	"io"
	"log"
	"music/utils"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func SearchVideo(name string) ([]utils.VideoSearch, error) {
	db := utils.StartConn()
	const ytUrl string = "https://music.youtube.com/youtubei/v1/search?prettyPrint=false"
	allVideos := make([]utils.VideoSearch, 0)

	var jsonStr = fmt.Sprintf(`{"context": {"client":{"clientName": "WEB_REMIX", "clientVersion": "1.20240827.03.00"}}, "params": "EgWKAQIIAWoQEAMQBBAJEAoQBRAREBAQFQ%%3D%%3D", "query": "%s"}`, name)

	// SELECT album.cover, album_music.music_id, music.id, music.title FROM album_music, music,album WHERE music.title LIKE '%Full%' AND album_music.music_id = music.id AND album.id = album_music.album_id;
	var sql string = "SELECT music.id, music.title FROM album_music JOIN music ON album_music.music_id = music.id JOIN album ON album.id = album_music.album_id WHERE music.title ~* $1;"
	rows, err := db.Query(sql, name)

	if err != nil {
		log.Println("Error querying db > ", err)
		return allVideos, nil
	}
	defer rows.Close()
	for rows.Next() {
		var musicListDb utils.MusicListDb
		err := rows.Scan(&musicListDb.Music_id, &musicListDb.Title)
		musicListDb.Cover = "/api/cover/" + musicListDb.Music_id

		if err != nil {
			log.Println("Error scanning rows > ", err)
			continue
		}
		videoData := &utils.VideoSearch{
			Id:       musicListDb.Music_id,
			Title:    musicListDb.Title,
			ImageUrl: musicListDb.Cover,
			Provider: "local",
		}
		allVideos = append(allVideos, *videoData)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error scanning rows > ", err)
		return allVideos, nil
	}

	if os.Getenv("INCLUDE_YOUTUBE") == "true" {
		res, err := http.Post(ytUrl, "application/json", strings.NewReader(jsonStr))
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		re := regexp.MustCompile(`"flexColumns":\[{"musicResponsiveListItemFlexColumnRenderer":\{"text":\{"runs":\[\{"text":"([^"]{0,100})".{0,600}"videoId":"([^"]{0,50})"`)
		matches := re.FindAllStringSubmatch(string(body), -1)

		for _, match := range matches {
			videoData := &utils.VideoSearch{
				Id:       match[2],
				Title:    match[1],
				ImageUrl: "https://i.ytimg.com/vi/" + match[2] + "/hqdefault.jpg",
				Provider: "youtube",
			}

			allVideos = append(allVideos, *videoData)
		}
	}

	return allVideos, nil

}
