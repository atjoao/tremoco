package functions

import (
	"io"
	"log"
	"music/server/env"
	"music/server/utils"
	"net/http"
	"net/url"
	"regexp"
)

func SearchVideo(name string) ([]utils.VideoSearch, error) {
	db := utils.StartConn()
	const ytUrl string = "https://www.youtube.com/results"
	allVideos := make([]utils.VideoSearch, 0)

	if env.INCLUDE_YOUTUBE {
		parseUrl, err := url.Parse(ytUrl)

		if err != nil {
			return nil, err
		}

		values := parseUrl.Query()
		values.Add("search_query", name)

		parseUrl.RawQuery = values.Encode()

		res, err := http.Get(parseUrl.String())
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// thx
		re := regexp.MustCompile(`"videoRenderer":\{"videoId":"([^"]{0,50})","thumbnail".{0,600}"title":\{"runs":\[\{"text":"([^"]{0,100})"`)
		matches := re.FindAllStringSubmatch(string(body), -1)

		for _, match := range matches {
			videoData := &utils.VideoSearch{
				Id:       match[1],
				Title:    match[2],
				ImageUrl: "https://i.ytimg.com/vi/" + match[1] + "/hqdefault.jpg",
				Provider: "youtube",
			}

			allVideos = append(allVideos, *videoData)
		}
	}

	// SELECT album.cover, album_music.music_id, music.id, music.title FROM album_music, music,album WHERE music.title LIKE '%Full%' AND album_music.music_id = music.id AND album.id = album_music.album_id;
	var sql string = "SELECT album.cover, music.id, music.title FROM album_music, music, album WHERE music.title ~* $1 AND album_music.music_id = music.id AND album.id = album_music.album_id;"
	rows, err := db.Query(sql, name)

	if err != nil {
		log.Println("Error querying db > ", err)
		return allVideos, nil
	}
	defer rows.Close()
	for rows.Next() {
		var musicListDb utils.MusicListDb
		err := rows.Scan(&musicListDb.Cover, &musicListDb.Music_id, &musicListDb.Title)
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

	return allVideos, nil

}
