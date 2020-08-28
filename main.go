package main

import (
	"os"
	"log"
	"time"
	"context"
	"strings"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

var cfg aws.Config

const layout string = "20060102150405.000"

func HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	t := time.Now()
	t_ := strings.Replace(t.Format(layout), ".", "", 1)
	if os.Getenv("RANDOM_VALUE") == t_[(len(t_) - len(os.Getenv("RANDOM_VALUE"))):] {
		client := cloudformation.New(cfg)
		req := client.DeleteStackRequest(&cloudformation.DeleteStackInput{
			StackName: aws.String(os.Getenv("STACK_NAME")),
		})
		_, err := req.Send(ctx)
		if err != nil {
			log.Println(err)
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            "<html><head><title>Serverless Application Ephemerality</title></head><body><span>Serverless Application Ephemerality</span></body></html>",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func init() {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
