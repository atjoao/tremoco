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
	VideoCodec   string `json:"videoCodec"`
}
