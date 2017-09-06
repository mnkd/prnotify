package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PullRequestUser struct {
	AvatarURL string `json:"avatar_url"`
	Login     string `json:"login"`
}

type PullRequest struct {
	Assignees          []PullRequestUser `json:"assignees"`
	RequestedReviewers []PullRequestUser `json:"requested_reviewers"`
	ClosedAt           interface{}       `json:"closed_at"`
	CreatedAt          string            `json:"created_at"`
	HTMLURL            string            `json:"html_url"`
	ID                 int64             `json:"id"`
	IssueURL           string            `json:"issue_url"`
	Number             int64             `json:"number"`
	State              string            `json:"state"`
	Title              string            `json:"title"`
	UpdatedAt          string            `json:"updated_at"`
	URL                string            `json:"url"`
	User               PullRequestUser   `json:"user"`
	Links              struct {
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
	} `json:"_links"`
}

func (gh GitHubAPI) GetPulls() ([]PullRequest, error) {
	var pulls []PullRequest

	// Prepare HTTP Request
	url := gh.BaseURL() + "/pulls" + "?access_token=" + gh.AccessToken
	req, err := http.NewRequest("GET", url, nil)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - PullRequest: <error> parse http request form:", parseFormErr)
		return pulls, parseFormErr
	}

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - PullRequest: <error> fetch pulls:", err)
		return pulls, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &pulls); err != nil {
		fmt.Fprintln(os.Stderr, "GitHubAPI - PullRequest: <error> json unmarshal:", err)
		return pulls, err
	}

	return pulls, nil
}
