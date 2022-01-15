package rapididentity

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func getCredentials() (string, string, error) {
	var host string
	var key string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return host, key, err
	}
	_, err = os.Stat(homeDir + credsFile)
	if err != nil {
		return host, key, err
	}

	credsData, err := os.Open(homeDir + credsFile)
	if err != nil {
		return host, key, err
	}

	credsDataScanner := bufio.NewScanner(credsData)
	for credsDataScanner.Scan() {
		if strings.Contains(credsDataScanner.Text(), "host") {
			host = strings.ReplaceAll(credsDataScanner.Text(), "host=", "")
		} else if strings.Contains(credsDataScanner.Text(), "service_key") {
			key = strings.ReplaceAll(credsDataScanner.Text(), "service_key=", "")
		}
	}

	return host, key, nil
}

func GetAnnouncements() (*[]Announcement, error) {

	host, key, err := getCredentials()

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []Announcement
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateAnnouncement(announcement *Announcement) (*Announcement, error) {

	host, key, err := getCredentials()

	if err != nil {
		return nil, err
	}

	announcement_obj, _ := json.Marshal(announcement)

	body := bytes.NewBuffer(announcement_obj)

	client := &http.Client{}
	req, err := http.NewRequest("POST", host+path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Announcement
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func DeleteAnnouncement(id string) error {

	host, key, err := getCredentials()

	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", host+path+"/"+id, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
