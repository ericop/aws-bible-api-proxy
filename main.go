package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	UrlText string `json:"urlText"`
	Code    string `json:"code"`
	// header x-api-key =
	// X-API-Key:eAamcrnwum9yI7J9lDPYp3zLnDrBoqLcaLKBDDjc
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(request Request) (Response, error) {
	if strings.Contains(request.UrlText, ".mp3") {
		log.Printf("Has MP3 in %v", request.UrlText)

		return Response{
			Message: fmt.Sprintf("You sent MP3 request urlText: `%v`, code: `%v`", request.UrlText, request.Code),
			Ok:      true,
		}, nil
		// http://fcbhabdm.s3.amazonaws.com/mp3audiobibles2/ENGESVO1DA/A19__002_Psalms______ENGESVO1DA.mp3

		downloadFile(fmt.Sprintf("%d.mp3", time.Now().Unix()), request.UrlText)
		//     var respMp3 = Task.Run(() => client.GetStreamAsync(urlText));
		//     var respAudioFile = respMp3.GetAwaiter().GetResult();
		//     var responseMp3 = new OkObjectResult(respAudioFile);

		//     responseMp3.ContentTypes.Add("audio/mpeg");
		//     return responseMp3 != null
		//          ? (ActionResult)responseMp3
		//          : new BadRequestObjectResult("Please pass a `url` on the query string or in the request body");

	}

	//apiKey := "5b50f7439b939d9f4faa4bf81e0c8f46"
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

	log.Printf("Has text not MP3 in %v", request.UrlText)
	return Response{
		Message: fmt.Sprintf("You sent urlText: `%v`, code: `%v`", request.UrlText, request.Code),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
