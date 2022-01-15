package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	snsEntity := snsEvent.Records[0].SNS
	switch snsEntity.Message {
	case "command":
		err := displayIdeaModal(&snsEntity)
		if err != nil {
			log.Fatal(err)
		}
	case "view":
		err := submitIdea(&ctx, &snsEntity, svc)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func displayIdeaModal(snsEntity *events.SNSEntity) error {

	attributes, err := getSNSAttributes(snsEntity.MessageAttributes)
	if err != nil {
		return err
	}
	body := ViewResponse{
		TriggerId: attributes["trigger_id"],
		View: SlackModal{
			Type: "modal",
			Title: SlackText{
				Type:  "plain_text",
				Text:  "Submit Idea",
				Emoji: true,
			},
			Submit: SlackText{
				Type:  "plain_text",
				Text:  "Submit",
				Emoji: true,
			},
			Close: SlackText{
				Type:  "plain_text",
				Text:  "Cancel",
				Emoji: true,
			},
			Blocks: []SlackBlock{
				{
					Type: "input",
					Element: SlackElement{
						Type: "static_select",
						Placeholder: &SlackText{
							Type:  "plain_text",
							Text:  "Select a market",
							Emoji: true,
						},
						Options: []SlackOption{
							{
								Text: SlackText{
									Type:  "plain_text",
									Text:  "K12",
									Emoji: true,
								},
								Value: "K12",
							},
							{
								Text: SlackText{
									Type:  "plain_text",
									Text:  "Higher Education",
									Emoji: true,
								},
								Value: "Higher Education",
							},
							{
								Text: SlackText{
									Type:  "plain_text",
									Text:  "K12 and Higher Education",
									Emoji: true,
								},
								Value: "K12 and Higher Education",
							},
						},
						ActionId: "static_select-action",
					},
					Label: SlackText{
						Type:  "plain_text",
						Text:  "Market",
						Emoji: true,
					},
				},
				{
					Type: "input",
					Element: SlackElement{
						Type:     "plain_text_input",
						ActionId: "plain_text_input-action",
					},
					Label: SlackText{
						Type:  "plain_text",
						Text:  "Idea Name",
						Emoji: true,
					},
				},
				{
					Type: "input",
					Element: SlackElement{
						Type:      "plain_text_input",
						Multiline: true,
						ActionId:  "plain_text_input-action",
					},
					Label: SlackText{
						Type:  "plain_text",
						Text:  "Describe your idea",
						Emoji: true,
					},
				},
				{
					Type: "input",
					Element: SlackElement{
						Type:      "plain_text_input",
						Multiline: true,
						ActionId:  "plain_text_input-action",
					},
					Label: SlackText{
						Type:  "plain_text",
						Text:  "How does this idea help the market selected",
						Emoji: true,
					},
				},
			},
		},
	}

	result, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", result)
	resultReader := bytes.NewReader(result)
	newRequest, err := http.NewRequest("POST", "https://slack.com/api/views.open", resultReader)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", newRequest)
	newRequest.Header.Set("Content-Type", "application/json")
	newRequest.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func submitIdea(ctx *context.Context, snsEntity *events.SNSEntity, svc *dynamodb.Client) error {
	fmt.Printf("%+v\n", snsEntity)
	attributes, err := getSNSAttributes(snsEntity.MessageAttributes)
	if err != nil {
		return err
	}
	today := time.Now().Format("2006-01-02T15:04:05")
	item := dynamodb.PutItemInput{
		TableName: aws.String("idea"),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: snsEntity.MessageID,
			},
			"targetMarket": &types.AttributeValueMemberS{
				Value: attributes["targetMarket"],
			},
			"ideaName": &types.AttributeValueMemberS{
				Value: attributes["ideaName"],
			},
			"ideaDescription": &types.AttributeValueMemberS{
				Value: attributes["ideaDescription"],
			},
			"valueDescription": &types.AttributeValueMemberS{
				Value: attributes["valueDescription"],
			},
			"createdAt": &types.AttributeValueMemberS{
				Value: today,
			},
			"updatedAt": &types.AttributeValueMemberS{
				Value: today,
			},
			"leanCanvas": &types.AttributeValueMemberBOOL{
				Value: false,
			},
			"completed": &types.AttributeValueMemberBOOL{
				Value: false,
			},
			"productAccepted": &types.AttributeValueMemberBOOL{
				Value: false,
			},
			"user": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"name": &types.AttributeValueMemberS{
						Value: attributes["displayName"],
					},
					"id": &types.AttributeValueMemberS{
						Value: attributes["userId"],
					},
				},
			},
		},
	}

	fmt.Printf("%+v\n", item)
	output, err := svc.PutItem(*ctx, &item)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", output)

	message := map[string]string{
		"channel": attributes["userId"],
		"text":    "First of all you are awesome! Thank you for submitting your idea this truly does help a ton! I will be reviewing your idea and will let you what the next steps are! Stay Tuned!",
	}

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resultReader := bytes.NewReader(result)
	newRequest, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", resultReader)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", newRequest)
	newRequest.Header.Set("Content-Type", "application/json")
	newRequest.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func getSNSAttributes(snsAttributes map[string]interface{}) (map[string]string, error) {
	attributes := make(map[string]string)
	for k := range snsAttributes {
		attr, ok := snsAttributes[k].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to parse SNS attribute as map[string]interface{}")
		}
		attrValue, ok := attr["Value"].(string)
		if !ok {
			return nil, fmt.Errorf("Unable to parse sns attribute value as a string")
		}
		attributes[k] = attrValue
	}

	return attributes, nil
}
