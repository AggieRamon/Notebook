package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
				Body:       "Unable to decode base64 body",
			}, err
		}
	} else {
		body = []byte(req.Body)
	}

	validRequest, err := isAuthorized(req.Headers["x-slack-request-timestamp"], string(body), req.Headers["x-slack-signature"])
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to validate Slack Request",
		}, err
	}

	if !validRequest {
		fmt.Println("Signatures are mismatched")
		return events.APIGatewayProxyResponse{StatusCode: 401}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to make aws connection",
		}, err
	}

	svc := sns.NewFromConfig(cfg)
	awsBody := bytes.NewReader(body)
	awsRequest, err := http.NewRequest("POST", "https://api.slack.com", awsBody)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to create new request to parse x-www-form-urlencoded body",
		}, err
	}

	awsRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err = awsRequest.ParseForm()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to parse x-www-form-urlencoded body",
		}, err
	}
	fmt.Printf("%+v\n", awsRequest.PostForm)
	var payload SlackRes
	err = json.Unmarshal([]byte(awsRequest.PostForm["payload"][0]), &payload)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error unmarshalling body",
		}, err
	}
	snsInput, err := innovateSubmission(&payload)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Unable to submit idea",
		}, err
	}

	fmt.Printf("%+v\n", snsInput)

	output, err := svc.Publish(ctx, snsInput)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error publishing to SNS Topic",
		}, err
	}

	fmt.Printf("%+v\n", output)

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func innovateSubmission(payload *SlackRes) (*sns.PublishInput, error) {
	fmt.Printf("%+v\n", payload.View.Blocks)
	defMap := map[int]string{
		0: "targetMarket",
		1: "ideaName",
		2: "ideaDescription",
		3: "valueDescription",
	}

	snsInput := sns.PublishInput{
		Message:  aws.String("view"),
		TopicArn: aws.String("SNS Topic ARN"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"userId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(payload.User.Id),
			},
			"displayName": {
				DataType:    aws.String("String"),
				StringValue: aws.String(payload.User.Name),
			},
			"triggerId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(payload.TriggerId),
			},
		},
	}

	for i := range payload.View.Blocks {
		var selectValue string
		stateValue := payload.View.State.Values[payload.View.Blocks[i].BlockId]
		if i == 0 {
			var selectAction SlackSelectAction
			data, err := json.Marshal(stateValue)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &selectAction)
			if err != nil {
				return nil, err
			}
			selectValue = selectAction.StaticSelectAction.SelectedOption.Value
		} else {
			var plainText SlackPlainTextInputAction
			data, err := json.Marshal(stateValue)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &plainText)
			if err != nil {
				return nil, err
			}
			selectValue = plainText.PlainTextInputAction.Value
		}
		attrKey := defMap[i]
		snsInput.MessageAttributes[attrKey] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(selectValue),
		}
	}

	return &snsInput, nil
}

func isAuthorized(timestamp string, body string, signature string) (bool, error) {
	hashString := fmt.Sprintf("v0:%s:%s", timestamp, body)
	hash := hmac.New(sha256.New, []byte(os.Getenv("SECRET")))
	_, err := hash.Write([]byte(hashString))
	if err != nil {
		return false, nil
	}
	computedHash := fmt.Sprintf("v0=%s", hex.EncodeToString(hash.Sum(nil)))
	fmt.Printf("Computed: %s\nOriginal: %s\n", computedHash, signature)
	if computedHash == signature {
		return true, nil
	} else {
		return false, nil
	}
}
