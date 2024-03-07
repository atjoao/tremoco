package utils

type User struct {
	Username  string
	Password  string
	Playlists []Playlist
}

type Playlist struct {
	PlaylistId   int
	PlaylistName string
	MusicList    []VideoMeta
}

type VideoSearch struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"thumbnail"`
	Provider string `json:"provider"`
}

type VideoMeta struct {
	VideoId    string      `json:"videoid"`
	Title      string      `json:"title"`
	Author     string      `json:"author"`
	Duration   string      `json:"duration"`
	Thumbnails []Thumbnail `json:"thumbnails"`
	Streams    []Streams   `json:"streams"`
}

type Streams struct {
	AudioQuality string `json:"audioQuality"`
	MimeType     string `json:"mimeType"`
	StreamUrl    string `json:"streamUrl"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type VideoPlaybackResponse struct {
	PlayabilityStatus struct {
		Status          string `json:"status"`
		PlayableInEmbed bool   `json:"playableInEmbed"`
	} `json:"playabilityStatus"`

	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		Formats          []struct {
			Itag            uint16 `json:"itag"`
			URL             string `json:"url"`
			MimeType        string `json:"mimeType"`
			Bitrate         uint32 `json:"bitrate"`
			Width           uint16 `json:"width"`
			Height          uint16 `json:"height"`
			LastModified    string `json:"lastModified"`
			Quality         string `json:"quality"`
			Xtags           string `json:"xtags"`
			FPS             uint8  `json:"fps"`
			QualityLabel    string `json:"qualityLabel"`
			ProjectionType  string `json:"projectionType"`
			AudioQuality    string `json:"audioQuality"`
			ApproxDuration  string `json:"approxDurationMs"`
			AudioSampleRate string `json:"audioSampleRate"`
			AudioChannels   uint8  `json:"audioChannels"`
		} `json:"formats"`

		AdaptiveFormats []struct {
			Itag      uint16 `json:"itag"`
			URL       string `json:"url"`
			MimeType  string `json:"mimeType"`
			Bitrate   uint32 `json:"bitrate"`
			Width     uint16 `json:"width"`
			Height    uint16 `json:"height"`
			InitRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"initRange"`
			IndexRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"indexRange"`
			LastModified   string `json:"lastModified"`
			ContentLength  string `json:"contentLength"`
			Quality        string `json:"quality"`
			FPS            uint8  `json:"fps"`
			QualityLabel   string `json:"qualityLabel"`
			ProjectionType string `json:"projectionType"`
			AudioQuality   string `json:"audioQuality"`
			AverageBitrate uint32 `json:"averageBitrate"`
			ColorInfo      struct {
				Primaries               string `json:"primaries"`
				TransferCharacteristics string `json:"transferCharacteristics"`
				MatrixCoefficients      string `json:"matrixCoefficients"`
			} `json:"colorInfo"`
			ApproxDuration string `json:"approxDurationMs"`
		} `json:"adaptiveFormats"`
	} `json:"streamingData"`

	VideoDetails struct {
		VideoId          string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelId        string   `json:"channelId"`
		IsOwnerViewing   bool     `json:"isOwnerViewing"`
		ShortDescription string   `json:"shortDescription"`
		IsCrawlable      bool     `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []Thumbnail `json:"thumbnails"`
		} `json:"thumbnail"`
		AllowRatings      bool   `json:"allowRatings"`
		ViewCount         string `json:"viewCount"`
		Author            string `json:"author"`
		IsPrivate         bool   `json:"isPrivate"`
		IsUnpluggedCorpus bool   `json:"isUnpluggedCorpus"`
		IsLiveContent     bool   `json:"isLiveContent"`
	} `json:"videoDetails"`
}

type FFProbeOutputResponse struct {
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Artist string `json:"ARTIST"`
			Title  string `json:"TITLE"`
			Album  string `json:"ALBUM"`
			Genre  string `json:"GENRE"`
		} `json:"tags"`
	} `json:"format"`
}

type MusicListDb struct {
	Cover    string
	Music_id string
	Id       string
	Title    string
}