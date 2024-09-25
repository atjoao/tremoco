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
		const ytUrl string = "https://www.youtube.com/youtubei/v1/player?key=AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w"
		var jsonStr = fmt.Sprintf(`{"videoId": "%s","context": {"client": {"clientName": "ANDROID","clientVersion": "17.36.4","androidSdkVersion": 31,"hl": "en","gl": "US","utcOffsetMinutes": 0}}}`, videoId)

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
