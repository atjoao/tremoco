package utils

import "io/fs"

var Assets fs.FS

func CreateAssets(embedded fs.FS) fs.FS {
	var err error
	Assets, err = fs.Sub(embedded, "assets")
	if err != nil {
		panic(err)
	}

	return Assets

}
