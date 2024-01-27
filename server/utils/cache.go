package utils

import (
	"log"
	"strconv"
	"time"
)

var cache = make(map[string]VideoPlaybackResponse)

func CreateCache(videoData VideoPlaybackResponse) bool {
	_, valueInCache := cache[videoData.VideoDetails.VideoId]

	if valueInCache{
		return false
	}

	cache[videoData.VideoDetails.VideoId] = videoData

	go func() {
		videoExpireTime, err := strconv.ParseInt(videoData.StreamingData.ExpiresInSeconds, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(time.Duration(videoExpireTime)*time.Second)

		delete(cache, videoData.VideoDetails.VideoId)
	}()

	return true
}

func GetFromCache(videoId string) (bool, *VideoPlaybackResponse){
	value, valueInCache := cache[videoId]
	if !valueInCache{
		return false, nil
	} else {
		return true, &value
	}

}