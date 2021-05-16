package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	UrlText string `json:"urlText"`
	Code    string `json:"code"`
	// header x-api-key =
	// X-API-Key:eAamcrnwum9yI7J9lDPYp3zLnDrBoqLcaLKBDDjc
}

type RawString string
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

// MarshalJSON returns *m as the JSON encoding of m.
func (m *RawString) MarshalJSON() ([]byte, error) {
	return []byte(*m), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawString) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("RawString: UnmarshalJSON on nil pointer")
	}
	*m += RawString(data)
	return nil
}

func Handler(request Request) (Response, error) {

	encodedUrlText := request.UrlText
	decodedUrlText, err := url.QueryUnescape(encodedUrlText)
	if err != nil {
		log.Fatal(err)
	}

	apiKey := "5b50f7439b939d9f4faa4bf81e0c8f46"
	urlToCall := fmt.Sprintf("%v&key=%v", decodedUrlText, apiKey)

	log.Printf("urlToCall:%v", urlToCall)

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
	// [{
	// 	"book_name": "Matthew",
	// 	"book_id": "Matt",
	// 	"book_order": "55",
	// 	"chapter_id": "2",
	// 	"chapter_title": "Chapter 2",
	// 	"verse_id": "13",
	// 	"verse_text": "Now when they had departed, behold, an angel of the Lord appeared to Joseph in a dream and said, “Rise, take the child and his mother, and flee to Egypt, and remain there until I tell you, for Herod is about to search for the child, to destroy him.” \n\t\t\t",
	// 	"paragraph_number": "2"
	//   }]
	if readErr != nil {
		log.Fatal(readErr)
	}

	var bibleVerses []BibleVerse

	json.Unmarshal([]byte(body), &bibleVerses)

	// string requestBody = await new StreamReader(req.Body).ReadToEndAsync();
	// dynamic data = JsonConvert.DeserializeObject(requestBody);
	// urlText = urlText ?? data?.urlText;
	// urlText = string.IsNullOrWhiteSpace(urlText) ? "http://dbt.io/library/volume?language_code=ENG&v=2" : urlText;
	// var urlToCall = $"{urlText}&key={apiKey}";
	// Console.WriteLine($"data was '{data}'");
	// Console.WriteLine($"urlText was '{urlText}'");
	// Console.WriteLine($"urlToCall was '{urlToCall}'");

	// var resp = Task.Run(() => client.GetStringAsync(urlToCall));
	// var respText = resp.GetAwaiter().GetResult();
	// var asJson = JsonConvert.DeserializeObject(respText);
	// var response = new OkObjectResult(asJson);

	// response.ContentTypes.Add("application/json");
	// return respText != null
	// 	 ? (ActionResult)response
	// 	 : new BadRequestObjectResult("Please pass a `url` on the query string or in the request body");

	log.Printf("body from bible Api:", body)
	return Response{
		Message: bibleVerses,
		//Message: fmt.Sprintf("You sent urlText: `%v`, code: `%v`", request.UrlText, request.Code),
		Ok: true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
