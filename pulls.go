package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GitHubAPI struct {
	AccessToken string
	Owner       string
	Repo        string
	UsersMap    UsersMap
}

type PullRequest struct {
	Assignees []struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"assignees"`
	Body      string      `json:"body"`
	ClosedAt  interface{} `json:"closed_at"`
	CreatedAt string      `json:"created_at"`
	DiffURL   string      `json:"diff_url"`
	HTMLURL   string      `json:"html_url"`
	ID        int64       `json:"id"`
	IssueURL  string      `json:"issue_url"`
	MergedAt  interface{} `json:"merged_at"`
	Milestone interface{} `json:"milestone"`
	Number    int64       `json:"number"`
	State     string      `json:"state"`
	Title     string      `json:"title"`
	UpdatedAt string      `json:"updated_at"`
	URL       string      `json:"url"`
	User      struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"user"`
}

func (pr PullRequest) SlackMessage(usersMap UsersMap) string {
	var str = fmt.Sprintf("\t<%s|#%d, %s>", pr.HTMLURL, pr.Number, pr.Title)

	if len(pr.Assignees) > 0 {
		str += " => "
	}

	for _, assignee := range pr.Assignees {
		name := assignee.Login
		if v, ok := usersMap[assignee.Login]; ok {
			name = v
		}
		str += "@" + name + " "
	}

	str += "\n"
	return str
}

func isWIP(title string) bool {
	return strings.Contains(strings.ToUpper(title), "WIP")
}

func (gh GitHubAPI) SlackMessage(pulls []PullRequest) string {
	if len(pulls) == 0 {
		return ""
	}

	var array []PullRequest
	for _, pull := range pulls {
		if isWIP(pull.Title) {
			continue
		}
		array = append(array, pull)
	}

	var str = fmt.Sprintf("I found %d open pull requests for %s/%s:\n", len(array), gh.Owner, gh.Repo)
	for _, pull := range array {
		str += pull.SlackMessage(gh.UsersMap)
	}
	return str
}

func (gh GitHubAPI) urlString() string {
	return "https://api.github.com/repos/" + gh.Owner + "/" + gh.Repo + "/pulls" + "?access_token=" + gh.AccessToken
}

func (gh GitHubAPI) GetPulls() ([]PullRequest, error) {
	var pulls []PullRequest

	// Prepare HTTP Request
	url := gh.urlString()
	req, err := http.NewRequest("GET", url, nil)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Fprintln(os.Stderr, "[Error] Parse HTTP request form:", parseFormErr)
		return pulls, parseFormErr
	}

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Fetch pulls: ", err)
		return pulls, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &pulls); err != nil {
		fmt.Fprintln(os.Stderr, "[Error] JSON unmarshal:", err)
		return pulls, err
	}

	return pulls, nil
}

func NewGitHubAPI() GitHubAPI {
	return GitHubAPI{}
}
