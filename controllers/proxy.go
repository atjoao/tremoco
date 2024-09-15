package controllers

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"

	"strings"

	"github.com/gin-gonic/gin"
)

var allowedDomains = []string{
	"i.ytimg.com",
	"googlevideo.com",
	"github.com",
	"raw.githubusercontent.com",
}

func isDomainAllowed(domain string, allowedDomains []string) bool {
	for _, allowedDomain := range allowedDomains {
		if strings.Contains(domain, allowedDomain) {
			return true
		}
	}
	return false
}

func ProxyContent(ctx *gin.Context) {
	urlQuery := ctx.Query("url")

	if len(urlQuery) <= 0 {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": `no "url" query param or empty`,
		})
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(urlQuery)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status":  "INVALID_PARAMS",
			"message": "Invalid base64 url",
		})
		return
	}

	urlParse, err := url.Parse(string(decoded))
	if err != nil || urlParse.Scheme == "" || urlParse.Host == "" || urlParse.Scheme != "http" && urlParse.Scheme != "https" {
		ctx.JSON(400, gin.H{
			"status":  "INVALID_PARAMS",
			"message": "Invalid url",
		})
		return
	}

	if !isDomainAllowed(urlParse.Host, allowedDomains) {
		ctx.JSON(400, gin.H{
			"status":  "INVALID_PARAMS",
			"message": "Invalid domain",
		})
		return
	}

	req, err := http.NewRequest("GET", string(decoded), nil)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error while creating request",
		})
		return
	}

	rangeHeader := ctx.GetHeader("Range")
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Error while fetching the content",
		})
		return
	}

	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "No content type header",
		})
		return
	}

	if strings.Contains(contentType, "image") || strings.Contains(contentType, "audio") {
		for k, v := range res.Header {
			for _, vv := range v {
				ctx.Writer.Header().Add(k, vv)
			}
		}

		ctx.Status(res.StatusCode)
		_, err = io.Copy(ctx.Writer, res.Body)
		if err != nil {
			ctx.JSON(500, gin.H{
				"status":  "SERVER_ERROR",
				"message": "Error while streaming content",
			})
		}
		return
	}

	ctx.JSON(400, gin.H{
		"status":  "INVALID_CONTENT",
		"message": "Invalid content type",
	})
}
