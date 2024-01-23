package functions

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ParseYTPlayer() {
	response, err := http.Get("https://www.youtube.com/iframe_api")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	scriptBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	script := string(scriptBytes)

	regexScrURL := regexp.MustCompile(`var scriptUrl = '([^']+)';`)
	match := regexScrURL.FindStringSubmatch(script)
	scriptURL := strings.Replace(match[1], "\\", "", -1)
	fmt.Println(scriptURL)

	regexPlayerID := regexp.MustCompile(`https:\/\/www\.youtube\.com\/s\/player\/([A-Za-z0-9]+)\/www-widgetapi\.vflset\/www-widgetapi\.js`)
	matchPlayer := regexPlayerID.FindStringSubmatch(scriptURL)
	fmt.Println(matchPlayer[1])

	playerURL := fmt.Sprintf("https://www.youtube.com/s/player/%s/player_ias.vflset/en_US/base.js", matchPlayer[1])
	fmt.Println(playerURL)
	playerScriptResponse, err := http.Get(playerURL)
	if err != nil {
		panic(err)
	}
	defer playerScriptResponse.Body.Close()

	playerScriptBytes, err := io.ReadAll(playerScriptResponse.Body)
	if err != nil {
		panic(err)
	}
	playerScript := string(playerScriptBytes)

	scriptRegex := regexp.MustCompile(`(?m)^(.*a=a\.split\(\"\"\);.*)$`)
	getFunction := scriptRegex.FindStringSubmatch(playerScript)
	fmt.Println("function dec:", getFunction[0])

	reFunctionName := regexp.MustCompile(`(?m)^.*a=a\.split\(""\);([^\.]{1,3}).*$`)
	getFunctionName := reFunctionName.FindStringSubmatch(playerScript)
	fmt.Println("function name:", getFunctionName[1])

	var reFunctionVarStr string = "var\\s+fn\\s*=\\s*{[^}]*}"
	replacedStr := strings.ReplaceAll(reFunctionVarStr, "fn", getFunctionName[1])
	reFunctionVar := regexp.MustCompile(replacedStr)
	getFunctionVar := reFunctionVar.FindString(playerScript)
	fmt.Println(getFunctionVar)

	reSignatureStamp := regexp.MustCompile(`signatureTimestamp:(\d+)`)
	getSignatureStamp := reSignatureStamp.FindStringSubmatch(playerScript)
	fmt.Println(getSignatureStamp[1])
}
