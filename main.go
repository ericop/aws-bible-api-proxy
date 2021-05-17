package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// type Event struct {
// 	IsBase64Encoded   bool    `json:"isBase64Encoded"`
// 	StatusCode        string  `json:"statusCode"`
// 	Headers           string  `json:"headers"`
// 	MultiValueHeaders string  `json:"multiValueHeaders"`
// 	Body              Request `json:"body"`
// }

type RequestBody struct {
	UrlText string `json:"urlText"`
	Code    string `json:"code"`
	// header X-API-Key: eAamcrnwum9yI7J9lDPYp3zLnDrBoqLcaLKBDDjc
	// header x-api-key: Genesis1-2InTheBeginningGodCreated
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

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("EVENT  =>`%v`", request)
	log.Printf("BODY  =>`%v`", request.Body)
	log.Printf("BODY as []byte =>`%v`", []byte(request.Body))
	log.Printf("HEADERS  =>`%v`", request.Headers)
	log.Printf("IsBase64Encoded  =>`%v`", request.IsBase64Encoded)

	var requestBody RequestBody

	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Fatal("Unmarshal request.Body failed", err)
	}

	encodedCode := requestBody.Code
	decodedCode, err := url.QueryUnescape(encodedCode)
	if err != nil {
		log.Fatal(err)
	}

	if decodedCode != "1nTh3B3g1nn1ngG0dCr34t3d" {
		return events.APIGatewayProxyResponse{Body: "invalid access code", StatusCode: 401}, nil
	}

	encodedUrlText := requestBody.UrlText
	decodedUrlText, err := url.QueryUnescape(encodedUrlText)
	if err != nil {
		log.Fatal(err)
	}

	isBibleTextRequest := strings.Contains(decodedUrlText, "/text/verse")
	isAudioLocationRequest := strings.Contains(decodedUrlText, "/audio/location")
	//isAudioPathRequest := strings.Contains(decodedUrlText, "/audio/path")

	apiKey := "5b50f7439b939d9f4faa4bf81e0c8f46"
	urlToCall := fmt.Sprintf("%v&key=%v", decodedUrlText, apiKey)
	//log.Printf("urlToCall:%v", urlToCall)

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
	//log.Printf("body from bible Api:", body)

	if isBibleTextRequest {
		var bibleVerses []BibleVerse
		json.Unmarshal([]byte(body), &bibleVerses)
		bibleVersesString, err := json.Marshal(bibleVerses)
		if err != nil {
			log.Fatalf("Error occurred during bibleVerses marshaling. Error: %s", err.Error())
		}
		return events.APIGatewayProxyResponse{Body: string(bibleVersesString), StatusCode: 200}, nil
	}

	if isAudioLocationRequest {
		var audioLocations []BibleAudioLocation
		json.Unmarshal([]byte(body), &audioLocations)

		audioLocationsString, err := json.Marshal(audioLocations)
		if err != nil {
			log.Fatalf("Error occurred during audioLocations marshaling. Error: %s", err.Error())
		}
		return events.APIGatewayProxyResponse{Body: string(audioLocationsString), StatusCode: 200}, nil
	}

	var audioPaths []BibleAudioPath
	json.Unmarshal([]byte(body), &audioPaths)

	audioPathsString, err := json.Marshal(audioPaths)
	if err != nil {
		log.Fatalf("Error occurred during audioPaths marshaling. Error: %s", err.Error())
	}
	return events.APIGatewayProxyResponse{Body: string(audioPathsString), StatusCode: 200}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
