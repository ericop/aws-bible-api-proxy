package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type RequestBody struct {
	UrlText string `json:"urlText"`
}

type Response struct {
	Message []BibleVerse `json:"message"`
	Ok      bool         `json:"ok"`
}

type BibleVerse struct {
	Book_name        string `json:"book_name"`
	Book_id          string `json:"book_id"`
	Book_order       string `json:"book_order"`
	Chapter_id       string `json:"chapter_id"`
	Chapter_title    string `json:"chapter_title"`
	Verse_id         string `json:"verse_id"`
	Verse_text       string `json:"verse_text"`
	Paragraph_number string `json:"paragraph_number"`
}

type BibleAudioLocation struct {
	Server    string `json:"server"`
	Root_path string `json:"root_path"`
	Protocol  string `json:"protocol"`
	CDN       string `json:"CDN"`
	Priority  string `json:"priority"`
}

type BibleAudioPath struct {
	Book_id    string `json:"book_id"`
	Chapter_id string `json:"chapter_id"`
	Path       string `json:"path"`
}

func Handler(requestBody RequestBody) (interface{}, error) {
	//log.Printf("typesIn requestBody  =>`%v`", requestBody)

	encodedUrlText := requestBody.UrlText
	decodedUrlText, err := url.QueryUnescape(encodedUrlText)
	if err != nil {
		log.Fatal(err)
	}

	isBibleTextRequest := strings.Contains(decodedUrlText, "/text/verse")
	isAudioLocationRequest := strings.Contains(decodedUrlText, "/audio/location")

	thirdPartyBibleApiKey := "5b50f7439b939d9f4faa4bf81e0c8f46"
	urlToCall := fmt.Sprintf("%v&key=%v", decodedUrlText, thirdPartyBibleApiKey)

	bibleClient := http.Client{
		Timeout: time.Second * 15, // timeout after 15 seconds
	}

	req, err := http.NewRequest(http.MethodGet, urlToCall, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, getErr := bibleClient.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	if isBibleTextRequest {
		var bibleVerses []BibleVerse
		json.Unmarshal([]byte(body), &bibleVerses)

		return bibleVerses, nil
	}

	if isAudioLocationRequest {
		var audioLocations []BibleAudioLocation
		json.Unmarshal([]byte(body), &audioLocations)

		return audioLocations, nil
	}

	// isAudioPathRequest := strings.Contains(decodedUrlText, "/audio/path")
	var audioPaths []BibleAudioPath
	json.Unmarshal([]byte(body), &audioPaths)

	return audioPaths, nil

}

func main() {
	lambda.Start(Handler)
}
