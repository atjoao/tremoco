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

func VideoMeta(videoId string, includeVideo bool) (*utils.VideoPlaybackResponse, []utils.Streams, error) {
    metas := make([]utils.Streams, 0)
    var response utils.VideoPlaybackResponse

    inCache, getCacheValue := utils.GetFromCache(videoId)
    if !inCache{
        const ytUrl string = "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
        var jsonStr = fmt.Sprintf(`{"videoId": "%s","context": {"client": {"clientName": "ANDROID_TESTSUITE","clientVersion": "1.9","androidSdkVersion": 30,"hl": "en","gl": "US","utcOffsetMinutes": 0}}}`, videoId)

        getVideoInfo, err := http.Post(ytUrl, "application/json", strings.NewReader(jsonStr))
        if err != nil {
            return nil, nil, err
        }

        defer getVideoInfo.Body.Close()

        getVideoBody, err := io.ReadAll(getVideoInfo.Body)
        if err != nil {
            return nil, nil, err
        }

        err = json.Unmarshal(getVideoBody, &response)
        utils.CreateCache(response)
        if err != nil {
            return nil, nil, err
        }
    } else {
        response = *getCacheValue
    }

	for _, data := range response.StreamingData.AdaptiveFormats {
		encodedUrl, err := url.QueryUnescape(data.URL)
		if err != nil {
			return nil, nil, err
		}

		if !includeVideo && data.AudioQuality == "" {
			continue
		}

		videoMeta := &utils.Streams{
			AudioQuality: data.AudioQuality,
			StreamUrl:    encodedUrl,
			MimeType:     data.MimeType,
		}

		metas = append(metas, *videoMeta)
	}

	return &response, metas, nil
}
