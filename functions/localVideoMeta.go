package functions

import (
	"music/server/utils"
)

func LocalVideoMeta(videoId string) (*utils.VideoMeta ,string){
	db := utils.StartConn()
	var err error

	var sql string = "SELECT album.cover, music.id, music.title, music.duration, music.author, music.location FROM album_music, music, album WHERE music.id = $1 AND album_music.music_id = music.id AND album.id = album_music.album_id"
	var music utils.VideoMeta
	var location string
	music.Thumbnails = append(music.Thumbnails, utils.Thumbnail{URL: ""})
	music.Streams = append(music.Streams, utils.Streams{AudioQuality: "", MimeType: "", StreamUrl: "/api/stream/" + videoId})
	
	err = db.QueryRow(sql, videoId).Scan(&music.Thumbnails[0].URL, &music.VideoId, &music.Title, &music.Duration, &music.Author, &location)
	if err != nil {
		return nil, ""
	}

	var output *utils.FFProbeOutputResponse
	output, err = FfprobeOutput(location)
	if err != nil {
		return nil, ""
	}
	music.Streams[0].MimeType = "audio/"+output.Format.FormatName
	music.Streams[0].AudioQuality = output.Format.Bitrate

	return &music, location
}