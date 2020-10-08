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
	"github.com/aws/aws-sdk-go-v2/config"
	slambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

var cfg aws.Config

func HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	count, _ := strconv.Atoi(os.Getenv("COUNT"))
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	if &request.RequestContext != nil && &request.RequestContext.HTTP != nil && len (request.RequestContext.HTTP.SourceIP) > 0 {
		log.Println(request.RequestContext.HTTP.SourceIP)
		goto Response
	}
	if count < limit {
		client := getLambdaClient()
		res, err := client.GetFunctionConfiguration(ctx, &slambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(os.Getenv("FUNCTION_NAME")),
		})
		if err != nil {
			log.Println(err)
		} else {
			env := res.Environment.Variables
			env["COUNT"] = aws.String(strconv.Itoa(count + 1))
			_, err := client.UpdateFunctionConfiguration(ctx, &slambda.UpdateFunctionConfigurationInput{
				FunctionName: aws.String(os.Getenv("FUNCTION_NAME")),
				Environment: &types.Environment{
					Variables: env,
				},
			})
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		client := getCloudformationClient()
		_, err := client.DeleteStack(ctx, &cloudformation.DeleteStackInput{
			StackName: aws.String(os.Getenv("STACK_NAME")),
		})
		if err != nil {
			log.Println(err)
		}
	}
Response:
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Body:            "<html><head><title>Serverless Application Ephemerality</title></head><body><span>Serverless Application Ephemerality</span></body></html>",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func getLambdaClient() *slambda.Client {
	if cfg.Region != os.Getenv("REGION") {
		cfg = getConfig()
	}
	return slambda.NewFromConfig(cfg)
}

func getCloudformationClient() *cloudformation.Client {
	if cfg.Region != os.Getenv("REGION") {
		cfg = getConfig()
	}
	return cloudformation.NewFromConfig(cfg)
}

func getConfig() aws.Config {
	var err error
	cfg.Region = os.Getenv("REGION")
	newConfig, err := config.LoadDefaultConfig()
	newConfig.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
	return newConfig
}

func main() {
	lambda.Start(HandleRequest)
}
