package utils

import (
	"log"
	"strconv"
	"time"
)

var streamCache = make(map[string]YT_VideoPlaybackResponse)

func StreamCreateCache(videoData YT_VideoPlaybackResponse) bool {
	_, valueInCache := streamCache[videoData.VideoDetails.VideoId]

	if valueInCache{
		return false
	}

	streamCache[videoData.VideoDetails.VideoId] = videoData

	go func() {
		videoExpireTime, err := strconv.ParseInt(videoData.StreamingData.ExpiresInSeconds, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(time.Duration(videoExpireTime)*time.Second)

		delete(streamCache, videoData.VideoDetails.VideoId)
	}()

	return true
}

func StreamGetFromCache(videoId string) (bool, *YT_VideoPlaybackResponse){
	value, valueInCache := streamCache[videoId]
	if !valueInCache{
		return false, nil
	} else {
		return true, &value
	}

}