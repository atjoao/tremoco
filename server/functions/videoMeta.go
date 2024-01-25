package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"music/server/utils"
	"net/http"
	"net/url"
	"strings"
)

func VideoMeta(videoId string, includeVideo bool) ([]utils.VideoMeta, error) {
	metas := make([]utils.VideoMeta, 0)
	const ytUrl string = "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
	var jsonStr = fmt.Sprintf(`{"videoId": "%s","context": {"client": {"clientName": "ANDROID_TESTSUITE","clientVersion": "1.9","androidSdkVersion": 30,"hl": "en","gl": "US","utcOffsetMinutes": 0}}}`, videoId)
	getVideoInfo, err := http.Post(ytUrl, "application/json", strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	defer getVideoInfo.Body.Close()

	getVideoBody, err := io.ReadAll(getVideoInfo.Body)
	if err != nil{
		return nil, err
	}

	var response utils.VideoPlaybackResponse 
	err = json.Unmarshal(getVideoBody, &response)
	if err != nil {
		return nil, err
	}

	for _, data := range response.StreamingData.AdaptiveFormats {
		encodedUrl, err := url.QueryUnescape(data.URL)
		if err != nil {
			fmt.Println(err)
		}

		if !includeVideo && data.AudioQuality == ""{
			continue
		}
		
		videoMeta := &utils.VideoMeta{
			AudioQuality: data.AudioQuality,
			StreamUrl:    encodedUrl,
			MimeType:     data.MimeType,
		}

		metas = append(metas, *videoMeta)
	}

	return metas, nil
}
