package utils

import (
	"log"
	"strconv"
	"time"
)

var streamCache = make(map[string]YT_VideoPlaybackResponse)
var localCache = make(map[string]VideoMeta)

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

func LocalCreateCache(music VideoMeta) bool {
	_, valueInCache := localCache[music.VideoId]

	if valueInCache{
		return false
	}

	localCache[music.VideoId] = music

	return true
}

func LocalGetFromCache(videoId string) (bool, *VideoMeta){
	value, valueInCache := localCache[videoId]
	if !valueInCache{
		return false, nil
	} else {
		return true, &value
	}
}

func StreamGetFromCache(videoId string) (bool, *YT_VideoPlaybackResponse){
	value, valueInCache := streamCache[videoId]
	if !valueInCache{
		return false, nil
	} else {
		return true, &value
	}

}