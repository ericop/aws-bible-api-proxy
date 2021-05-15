package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Name       string `json:"name"`
	IsLearning bool   `json:"isLearning"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(request Request) (Response, error) {
	learningPhrase := "start"
	if request.IsLearning {
		learningPhrase = "keep"
	}
	return Response{
		Message: fmt.Sprintf("Howdy %v! And %v learning. ;)", request.Name, learningPhrase),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
