package functions

import (
	"fmt"
	"io"
	"music/server/utils"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func VideoMeta(videoId string) ([]utils.VideoMeta, error) {
	metas := make([]utils.VideoMeta, 0)
	const ytUrl string = "https://www.youtube.com/watch"

	parseUrl, err := url.Parse(ytUrl)

	if err != nil {
		return nil, err
	}

	values := parseUrl.Query()
	values.Add("v", videoId)

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

	//re := regexp.MustCompile(`"audioQuality":"([^"]*)".*?"url":"([^"]*)".*?"mimeType":"audio/webm; codecs=\\\"([^\\\"]*)`)
	//re := regexp.MustCompile(`"audioQuality":"AUDIO_QUALITY_MEDIUM".*?"url":"([^"]*)".*?"mimeType":"audio/webm; codecs=\\\"([^\\\"]*)`)
	re := regexp.MustCompile(`"audioQuality":"([^"]*)".*?"url":"([^"]*)".*?"mimeType":"(audio|video)/webm; codecs=\\\"([^\\\"]*)`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		encodedUrl, err := url.QueryUnescape(match[2])
		if err != nil {
			fmt.Println(err)
		}

		encodedUrl = strings.Replace(encodedUrl, "\\u0026", "&", -1)

		videoMeta := &utils.VideoMeta{
			AudioQuality: match[1],
			StreamUrl:    encodedUrl,
			MimeType:     match[3],
			VideoCodec:   match[4],
		}

		metas = append(metas, *videoMeta)

	}
	return metas, nil
}
