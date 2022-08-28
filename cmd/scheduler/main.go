package main

import (
	"context"
	"jira-hermes/internal/app/server"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) {
	server.NewServer()
}

func main() {
	lambda.Start(HandleRequest)
}
