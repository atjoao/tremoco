package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func main() {
	
	res, err := http.Get("https://youtube.com/watch?v=b8orB2dMUKQ")
	if err != nil{	
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil{
		fmt.Println(err)

	}
	//re := regexp.MustCompile(`"audioQuality":"([^"]*)".*?"url":"([^"]*)".*?"mimeType":"audio/webm; codecs=\\\"([^\\\"]*)`)
	//re := regexp.MustCompile(`"audioQuality":"AUDIO_QUALITY_MEDIUM".*?"url":"([^"]*)".*?"mimeType":"audio/webm; codecs=\\\"([^\\\"]*)`)
	re := regexp.MustCompile(`"audioQuality":"([^"]*)".*?"url":"([^"]*)".*?"mimeType":"(audio|video)/webm; codecs=\\\"([^\\\"]*)`)

	matches := re.FindAllStringSubmatch(string(body), -1)

	for _, match := range matches {
		fmt.Println("---")

		fmt.Println("quality", match[1])
		fmt.Println("url - encoded", match[2])
		encodedUrl, err:= url.QueryUnescape(match[2])
		if err != nil{
			fmt.Println(err)
		}
		
		
		fmt.Println(strings.Replace(encodedUrl, "\\u0026","&", -1 ))

		fmt.Println("mimetype", match[3])
		fmt.Println("codec", match[4])
	}

}