package main

// Modify issues to report the results in age categories, say less than a month old, less than a year old, and more than a year old

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

const aMonth = 31 * 24
const aYear = 24 * 365

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues: \n", result.TotalCount)
	var inMonth, inYear, outYear []*github.Issue
	for _, item := range result.Items {
		creationDiff := hoursAgo(item.CreatedAt)
		fmt.Println(creationDiff, item.CreatedAt)
		// break
		if creationDiff < aMonth {
			inMonth = append(inMonth, item)
		} else if creationDiff < aYear {
			inYear = append(inYear, item)
		} else {
			outYear = append(outYear, item)
		}
	}
	fmt.Println("Issues added less than a month")
	for _, item := range inMonth {
		fmt.Printf("#%-5d %9.9s %v %.55s\n", item.Number, item.User.Login, item.CreatedAt, item.Title)
	}

	fmt.Println("Issues added less than a year more than a month")
	for _, item := range inYear {
		fmt.Printf("#%-5d %9.9s %v %.55s\n", item.Number, item.User.Login, item.CreatedAt, item.Title)
	}

	fmt.Println("Issues added more than year ago")
	for _, item := range outYear {
		fmt.Printf("#%-5d %9.9s %v %.55s\n", item.Number, item.User.Login, item.CreatedAt, item.Title)
	}
}

func hoursAgo(t time.Time) int {
	return int(time.Since(t).Hours())
}
