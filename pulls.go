package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/m-nakada/slackposter"
)

type PullRequest struct {
	Assignees []struct {
		Login string `json:"login"`
	} `json:"assignees"`
	ClosedAt  interface{} `json:"closed_at"`
	CreatedAt string      `json:"created_at"`
	HTMLURL   string      `json:"html_url"`
	ID        int64       `json:"id"`
	IssueURL  string      `json:"issue_url"`
	Number    int64       `json:"number"`
	State     string      `json:"state"`
	Title     string      `json:"title"`
	UpdatedAt string      `json:"updated_at"`
	URL       string      `json:"url"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
	Links struct {
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
	} `json:"_links"`
}

func (pr PullRequest) SlackAttachment(usersMap UsersMap) slackposter.Attachment {
	attachment := slackposter.Attachment{
		Fallback: fmt.Sprintf("#%d, %s (%s)", pr.Number, pr.Title, pr.HTMLURL, pr.assigneesString(usersMap, true)),
		Text:     pr.titleString() + " " + pr.assigneesString(usersMap, true),
		Color:    "warning",
	}

	return attachment
}

func (pr PullRequest) titleString() string {
	return fmt.Sprintf("\t<%s|#%d, %s> by @%s", pr.HTMLURL, pr.Number, pr.Title, pr.User.Login)
}

func (pr PullRequest) assigneesString(usersMap UsersMap, needsArrow bool) string {
	var str = ""
	if needsArrow && len(pr.Assignees) > 0 {
		str += " => "
	}
	for _, assignee := range pr.Assignees {
		name := assignee.Login
		if v, ok := usersMap[assignee.Login]; ok {
			name = v
		}
		str += "@" + name + " "
	}

	return str
}

func (pr PullRequest) SlackMessage(usersMap UsersMap) string {
	var str = pr.titleString()
	str += pr.assigneesString(usersMap, true)
	str += "\n"
	return str
}

func isWIP(title string) bool {
	return strings.Contains(strings.ToUpper(title), "WIP")
}

func filterdPulls(pulls []PullRequest) []PullRequest {
	var array []PullRequest
	for _, pull := range pulls {
		if isWIP(pull.Title) {
			continue
		}
		array = append(array, pull)
	}
	return array
}

func summaryString(owner string, repo string, pullsCount int) string {
	url := "https://github.com/" + owner + "/" + repo
	link := fmt.Sprintf("<%s|%s/%s>", url, owner, repo)

	var summary string
	switch pullsCount {
	case 0:
		summary = fmt.Sprintf("There's no open pull request for %s :tada: Let's take a break :dango: :tea:", link)
	case 1:
		summary = fmt.Sprintf("There's only one open pull request for %s :point_up:", link)
	default:
		summary = fmt.Sprintf("I found %d open pull requests for %s:\n", pullsCount, link)
	}
	return summary
}

func (gh GitHubAPI) SlackPayload(pulls []PullRequest) slackposter.Payload {
	var payload slackposter.Payload

	if len(pulls) == 0 {
		return payload
	}

	filterd := filterdPulls(pulls)
	var attachments []slackposter.Attachment

	for _, pull := range filterd {
		attachment := pull.SlackAttachment(gh.UsersMap)
		attachments = append(attachments, attachment)
	}

	payload.Text = summaryString(gh.Owner, gh.Repo, len(filterd))
	payload.Attachments = attachments

	return payload
}

func (gh GitHubAPI) SlackMessage(pulls []PullRequest) string {
	if len(pulls) == 0 {
		return ""
	}

	filterd := filterdPulls(pulls)
	var str = summaryString(gh.Owner, gh.Repo, len(filterd))
	for _, pull := range filterd {
		str += pull.SlackMessage(gh.UsersMap)
	}
	return str
}

func (gh GitHubAPI) GetPulls() ([]PullRequest, error) {
	var pulls []PullRequest

	// Prepare HTTP Request
	url := "https://api.github.com/repos/" + gh.Owner + "/" + gh.Repo + "/pulls" + "?access_token=" + gh.AccessToken
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
