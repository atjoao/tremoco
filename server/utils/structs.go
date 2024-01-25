package utils

type VideoSearch struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"thumbnail"`
}

type VideoMeta struct {
	AudioQuality string `json:"audioQuality"`
	MimeType     string `json:"mimeType"`
	StreamUrl    string `json:"streamUrl"`
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
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		AllowRatings      bool   `json:"allowRatings"`
		ViewCount         string `json:"viewCount"`
		Author            string `json:"author"`
		IsPrivate         bool   `json:"isPrivate"`
		IsUnpluggedCorpus bool   `json:"isUnpluggedCorpus"`
		IsLiveContent     bool   `json:"isLiveContent"`
	} `json:"videoDetails"`
}