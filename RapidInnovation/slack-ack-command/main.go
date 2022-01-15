package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body []byte
	var err error

	if req.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			fmt.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error decoding body from base 64 to string",
			}, err
		}
	} else {
		body = []byte(req.Body)
	}

	fmt.Printf("%+v\n", req)
	validRequest, err := isAuthorized(req.Headers["x-slack-request-timestamp"], string(body), req.Headers["x-slack-signature"])
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to validate Slack Signature",
		}, err
	}

	if !validRequest {
		return events.APIGatewayProxyResponse{StatusCode: 401}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error making AWS Connection",
		}, err
	}

	svc := sns.NewFromConfig(cfg)
	awsBody := bytes.NewReader(body)
	awsRequest, err := http.NewRequest("POST", "https://api.slack.com", awsBody)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating request to parse application/x-www-form-urlencoded body",
		}, err
	}
	awsRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = awsRequest.ParseForm()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to parse body",
		}, err
	}
	fmt.Printf("%+v\n", awsRequest.PostForm)
	snsInput := sns.PublishInput{
		Message:  aws.String("command"),
		TopicArn: aws.String("SNS Topic ARN"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"trigger_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(awsRequest.PostForm["trigger_id"][0]),
			},
		},
	}
	output, err := svc.Publish(ctx, &snsInput)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to publish to SNS Topic",
		}, err
	}
	fmt.Printf("%+v\n", output)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func isAuthorized(timestamp string, body string, signature string) (bool, error) {
	hashString := fmt.Sprintf("v0:%s:%s", timestamp, body)
	hash := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
	_, err := hash.Write([]byte(hashString))
	if err != nil {
		return false, err
	}
	computedHash := fmt.Sprintf("v0=%s", hex.EncodeToString(hash.Sum(nil)))
	fmt.Printf("Computed: %s\nOriginal: %s\n", computedHash, signature)
	if computedHash == signature {
		return true, nil
	} else {
		return false, nil
	}
}
