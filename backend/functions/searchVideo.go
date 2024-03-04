package functions

import (
	"io"
	structs "music/server/utils"
	"net/http"
	"net/url"
	"regexp"
)

func SearchVideo(name string) ([]structs.VideoSearch, error) {
	const ytUrl string = "https://www.youtube.com/results"
	allVideos := make([]structs.VideoSearch, 0)

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
		videoData := &structs.VideoSearch{
			Id:       match[1],
			Title:    match[2],
			ImageUrl: "https://i.ytimg.com/vi/" + match[1] + "/hqdefault.jpg",
		}

		allVideos = append(allVideos, *videoData)
	}

	// remember to append database videos

	return allVideos, nil

}
