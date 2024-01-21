package functions

import (
	"io"
	"music/server/structs"
	"net/http"
	"net/url"
	"regexp"
)



func SearchVideo(name string) ([]structs.Video, error){ 
	allVideos := make([]structs.Video, 0)
	ytUrl := "https://www.youtube.com/results"
	parseUrl, err := url.Parse(ytUrl)

	if err != nil{
		return nil, err
	}

	values := parseUrl.Query()
	values.Add("search_query", name)

	parseUrl.RawQuery = values.Encode()

	res, err := http.Get(parseUrl.String())
	if err != nil{
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil{
		return nil, err

	}

	// thx
	re := regexp.MustCompile(`"videoRenderer":\{"videoId":"([^"]{0,50})","thumbnail".{0,600}"title":\{"runs":\[\{"text":"([^"]{0,100})"`)
	matches := re.FindAllStringSubmatch(string(body), -1)


	for _, match := range matches {
		videoData := &structs.Video{
			Id:       match[1],
			Title:     match[2],
			ImageUrl: "https://i.ytimg.com/vi/"+ match[1] +"/hq720.jpg",
		}

		allVideos = append(allVideos, *videoData)
	}

	return allVideos, nil

}