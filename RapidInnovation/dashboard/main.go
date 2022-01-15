package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Idea struct {
	Id               string
	LeanCanvas       bool
	Completed        bool
	ProductAccepted  bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	IdeaDescription  string
	IdeaName         string
	User             UserInfo
	ValueDescription string
	TargetMarket     string
}

type UserInfo struct {
	Name string
	Id   string
}

func main() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	ideaTable := dynamodb.ScanInput{
		TableName: aws.String("idea"),
	}

	output, err := svc.Scan(ctx, &ideaTable)
	if err != nil {
		log.Fatal(err)
	}
	var ideas []Idea
	for _, v := range output.Items {
		itemMap := make(map[string]interface{})
		var idea Idea
		for key, val := range v {
			switch item := val.(type) {
			case *types.AttributeValueMemberS:
				if key == "createdAt" || key == "updatedAt" {
					d, _ := time.Parse("2006-01-02T15:04:05", item.Value)
					itemMap[strings.Title(key)] = d
				} else {
					itemMap[strings.Title(key)] = item.Value
				}
			case *types.AttributeValueMemberBOOL:
				itemMap[strings.Title(key)] = item.Value
			case *types.AttributeValueMemberM:
				userMap := make(map[string]string)
				userMap["Name"] = item.Value["name"].(*types.AttributeValueMemberS).Value
				userMap["Id"] = item.Value["id"].(*types.AttributeValueMemberS).Value
				itemMap["User"] = userMap
			}
		}
		result, err := json.Marshal(itemMap)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(result, &idea)
		if err != nil {
			log.Fatal(err)
		}
		ideas = append(ideas, idea)
	}
	sort.Slice(ideas, func(i, j int) bool { return ideas[i].CreatedAt.After(ideas[j].CreatedAt) })

	for _, idea := range ideas {
		loc, _ := time.LoadLocation("America/Chicago")
		d := idea.CreatedAt.In(loc)
		fmt.Printf("Total Number of Ideas: %v\n", len(output.Items))
		fmt.Printf("Market: %s\nIdea Name: %s\nIdea Description: %s\nValue Description: %s\nUser: %s\nCreated: %s\n----------------------------------------------\n", idea.TargetMarket, idea.IdeaName, idea.IdeaDescription, idea.ValueDescription, idea.User.Name, d)
	}
}
