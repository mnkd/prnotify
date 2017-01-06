package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Comment struct {
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	HTMLURL   string `json:"html_url"`
	ID        int64  `json:"id"`
	IssueURL  string `json:"issue_url"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (gh GitHubAPI) getCommentsWithURL(url string) ([]Comment, error) {
	var comments []Comment
	req, err := http.NewRequest("GET", url, nil)

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Fprintln(os.Stderr, "[Error] Parse HTTP request form:", parseFormErr)
		return comments, parseFormErr
	}

	// Fetch Request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Fetch comments: ", err, url)
		return comments, err
	}

	// Read Response Body
	resBody, _ := ioutil.ReadAll(res.Body)

	// Decode JSON
	if err := json.Unmarshal(resBody, &comments); err != nil {
		fmt.Fprintln(os.Stderr, "[Error] JSON unmarshal:", err)
		return comments, err
	}

	return comments, nil
}

func (gh GitHubAPI) GetCommentsWithPullRequest(pr PullRequest) ([]Comment, error) {
	var comments []Comment
	query := fmt.Sprintf("?access_token=%s&per_page=%d", gh.AccessToken, gh.Comment.PerPage)

	items, err := gh.getCommentsWithURL(pr.Links.Comments.Href + query)
	if err != nil {
		return comments, err
	} else {
		comments = append(comments, items...)
	}

	items, err = gh.getCommentsWithURL(pr.Links.ReviewComments.Href + query)
	if err == nil {
		return comments, err
	} else {
		comments = append(comments, items...)
	}

	return comments, nil
}
