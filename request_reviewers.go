package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type RequestedReviewer struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

func (gh GitHubAPI) GetRequestedReviewers(pull PullRequest) ([]RequestedReviewer, error) {
	var reviewers []RequestedReviewer

	// Prepare HTTP Request
	num := fmt.Sprintf("%d", pull.Number)
	url := gh.BaseURL() + "/pulls/" + num + "/requested_reviewers" + "?access_token=" + gh.AccessToken
	req, err := http.NewRequest("GET", url, nil)

	// Headers
	req.Header.Add("Accept", "application/vnd.github.black-cat-preview+json")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - Requested Reviewers: <error> parse http request form:", parseFormErr)
		return reviewers, parseFormErr
	}

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - Requested Reviewers: <error> fetch reviewers:", err)
		return reviewers, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &reviewers); err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - Requested Reviewers: <error> json unmarshal:", err)
		return reviewers, err
	}

	return reviewers, nil
}
