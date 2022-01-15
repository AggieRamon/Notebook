package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"aggieramon.com/rapididentity"
)

type AnnouncementCommand struct {
	fs               *flag.FlagSet
	id               string
	startTime        string
	endTime          string
	filterAclEnabled bool
	filterAcl        string
	message          string
	showOnce         bool
}

func main() {
	ac := &AnnouncementCommand{
		fs: flag.NewFlagSet("Create", flag.ContinueOnError),
	}
	ac.fs.StringVar(&ac.id, "id", "", "ID of the announcement. Only needed for deletion of an announcement")
	ac.fs.StringVar(&ac.startTime, "startTime", "", "Start time of announcement in yyyy-MM-dd hhmm format")
	ac.fs.StringVar(&ac.endTime, "endTime", "", "End time of announcement in yyyy-MM-dd hhmm format")
	ac.fs.BoolVar(&ac.filterAclEnabled, "filterAclEnabled", false, "Enable attribute based access control")
	ac.fs.StringVar(&ac.filterAcl, "filterAcl", "", "LDAP Filter definition (&(department=Finance)(employeeType=Staff))")
	ac.fs.StringVar(&ac.message, "message", "", "Message of the announcement")
	ac.fs.BoolVar(&ac.showOnce, "showOnce", false, "Set to true to show message only once")

	if len(os.Args) < 2 {
		ac.fs.Usage()
		os.Exit(1)
	} else {
		ac.fs.Parse(os.Args[2:])
	}

	if os.Args[1] == "configure" {
		configure()
	} else if os.Args[1] == "ls" {
		lsAnnouncements()
	} else if os.Args[1] == "create" {
		if ac.endTime == "" || ac.message == "" {
			log.Fatal("End time and message are required to create an announcement")
		}

		acl := rapididentity.Acl{
			FilterAclEnabled: ac.filterAclEnabled,
			FilterAcl:        ac.filterAcl,
		}

		announcement := rapididentity.Announcement{
			StartTime: ac.startTime,
			EndTime:   ac.endTime,
			Acl:       &acl,
			Message:   ac.message,
			ShowOnce:  ac.showOnce,
		}
		create(&announcement)
	} else if os.Args[1] == "delete" {
		err := rapididentity.DeleteAnnouncement(ac.id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Announcement Deleted")
	} else if os.Args[1] == "help" {
		ac.fs.Usage()
	} else {
		ac.fs.Usage()
	}
}

func configure() {
	var host string
	var key string
	fmt.Print("Please enter your RapidIdentity URL: ")
	fmt.Scan(&host)
	fmt.Print("Please enter your service identity key: ")
	fmt.Scan(&key)
	data := fmt.Sprintf("[default]\nhost=%v\nservice_key=%v\n", host, key)
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(homeDir + "/.rapididentity")

	if err != nil {
		err = os.Mkdir(homeDir+"/.rapididentity", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ioutil.WriteFile(homeDir+"/.rapididentity/credential", []byte(data), 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func lsAnnouncements() {
	resp, err := rapididentity.GetAnnouncements()
	if err != nil {
		log.Fatal(err)
	}
	result, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result)
}

func create(announcement *rapididentity.Announcement) {
	var endDate time.Time
	var startDate time.Time

	endDate, err := time.Parse("2006-01-02 15:04", announcement.EndTime)
	if err != nil {
		log.Fatal("Unable to parse end time. Ensure the following format yyyy-MM-dd hh:mm")
	}

	if announcement.StartTime != "" {
		startDate, err = time.Parse("2006-01-02 15:04", announcement.StartTime)
		if err != nil {
			log.Fatal("Unable to parse end time. Ensure the following format yyyy-MM-dd hh:mm")
		}
		announcement.StartTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:00.000Z", startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute())
	}

	announcement.EndTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:00.000Z", endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute())
	resp, err := rapididentity.CreateAnnouncement(announcement)
	if err != nil {
		log.Fatal(err)
	}
	result, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result)
}
