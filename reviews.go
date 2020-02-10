package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Review struct {
	ID   int64 `json:"id"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Body string `json:"body"`

	// State: APPROVED, COMMENTED, REQUESTED_CHANGED
	State string `json:"state"`
}

func (gh GitHubAPI) GetReviewsWithPullRequest(pull PullRequest) ([]Review, error) {
	var reviews []Review

	// Prepare HTTP Request
	url := fmt.Sprintf("%s/pulls/%d/reviews", gh.BaseURL(), pull.Number)
	req, err := http.NewRequest("GET", url, nil)

	// Headers
	req.Header.Add("Accept", "application/vnd.github.black-cat-preview+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", gh.AccessToken))

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - Reviews: <error> fetch reviews:", err)
		return reviews, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &reviews); err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - Reviews: <error> json unmarshal:", err)
		return reviews, err
	}

	return reviews, nil
}

func (review Review) IsApproved() bool {
	return review.State == "APPROVED"
}

func (review Review) IsChangesRequested() bool {
	return review.State == "CHANGES_REQUESTED"
}
