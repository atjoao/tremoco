package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"tremoco/utils"
)

func VideoMeta(videoId string, includeVideo bool) (*utils.YT_VideoPlaybackResponse, []utils.Streams, error) {
	metas := make([]utils.Streams, 0)
	var response utils.YT_VideoPlaybackResponse

	inCache, getCacheValue := utils.StreamGetFromCache(videoId)
	if !inCache {
		const ytUrl string = "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
		var jsonStr = fmt.Sprintf(`{"contentCheckOk": true, "context": {"client": {"androidSdkVersion": 31,"clientName": "ANDROID","clientVersion": "17.36.4","gl": "US","hl": "en-GB","osName": "Android","osVersion": "12","platform": "MOBILE"},"user": {"lockedSafetyMode": false}},"racyCheckOk": true,"videoId": "%s"}`, videoId)

		getVideoInfo, err := http.NewRequest("POST", ytUrl, strings.NewReader(jsonStr))
		if err != nil {
			return nil, nil, err
		}

		getVideoInfo.Header.Add("content-type", "application/json")
		getVideoInfo.Header.Add("User-Agent", "com.google.android.youtube/17.36.4 (Linux; U; Android 12; GB) gzip")

		client := &http.Client{}
		resp, err := client.Do(getVideoInfo)

		if err != nil {
			fmt.Println("Error sending request:", err)
			return nil, nil, err
		}

		defer resp.Body.Close()

		getVideoBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal(getVideoBody, &response)
		utils.StreamCreateCache(response)
		if err != nil {
			return nil, nil, err
		}
	} else {
		response = *getCacheValue
	}

	for _, data := range response.StreamingData.AdaptiveFormats {
		/* encodedUrl, err := url.QueryUnescape(data.URL)
		if err != nil {
			return nil, nil, err
		} */

		if !includeVideo && data.AudioQuality == "" {
			continue
		}

		videoMeta := &utils.Streams{
			AudioQuality: data.AudioQuality,
			StreamUrl:    data.URL,
			MimeType:     data.MimeType,
		}

		metas = append(metas, *videoMeta)
	}

	return &response, metas, nil

}
