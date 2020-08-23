package main

import (
	"os"
	"log"
	"context"
	"strconv"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	slambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

var cfg aws.Config

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count, _ := strconv.Atoi(os.Getenv("COUNT"))
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	if count < limit {
		client := slambda.New(cfg)
		req := client.GetFunctionConfigurationRequest(&slambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(os.Getenv("FUNCTION_NAME")),
		})
		res, err := req.Send(ctx)
		if err != nil {
			log.Println(err)
		} else {
			env := res.GetFunctionConfigurationOutput.Environment.Variables
			env["COUNT"] = strconv.Itoa(count + 1)
			req_ := client.UpdateFunctionConfigurationRequest(&slambda.UpdateFunctionConfigurationInput{
				FunctionName: aws.String(os.Getenv("FUNCTION_NAME")),
				Environment: &slambda.Environment{
					Variables: env,
				},
			})
			_, err := req_.Send(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
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
