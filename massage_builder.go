package main

import (
	"fmt"
	"github.com/m-nakada/slackposter"
	"strings"
)

type MessageBuilder struct {
	GitHubOwner  string
	GitHubRepo   string
	UsersManager UsersManager
}

func (builder MessageBuilder) titleString(pull PullRequest) string {
	return fmt.Sprintf("\t<%s|#%d> %s by %s",
		pull.HTMLURL, pull.Number, pull.Title, pull.User.Login)
}

func (builder MessageBuilder) allAssigneeString(pull PullRequest) string {
	if len(pull.Assignees) == 0 {
		name := builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
		return "@" + name + " *Assignee の指定をお願いします*"
	}

	var str = ""
	for _, assignee := range pull.Assignees {
		name := builder.UsersManager.ConvertGitHubToSlack(assignee.Login)
		str += "@" + name + " "
	}
	return str
}

func (builder MessageBuilder) reviewerString(pull PullRequest, thumbsUppers []string) string {
	var str = ""
	for _, assignee := range pull.Assignees {
		assigneeLogin := builder.UsersManager.ConvertGitHubToSlack(assignee.Login)
		found := false
		for _, thumbsUpper := range thumbsUppers {
			if assigneeLogin == thumbsUpper {
				found = true
				break
			}
		}
		if found == false {
			str += "@" + assigneeLogin + " "
		}
	}
	return str
}

func (builder MessageBuilder) BudildSummary(pullsCount int) string {
	repo := builder.GitHubOwner + "/" + builder.GitHubRepo
	url := "https://github.com/" + repo
	link := fmt.Sprintf("<%s|%s>", url, repo)

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

func (builder MessageBuilder) BuildAttachment(pull PullRequest, comments []Comment) slackposter.Attachment {

	var thumbsUppers []string

	for _, comment := range comments {
		if strings.Contains(comment.Body, ":+1:") || strings.Contains(comment.Body, "👍") {
			username := builder.UsersManager.ConvertGitHubToSlack(comment.User.Login)
			thumbsUppers = append(thumbsUppers, username)
		}
	}

	title := builder.titleString(pull)
	var color, reaction, mention string
	switch len(thumbsUppers) {
	case 0:
		reaction = ""
		color = "danger"
		name := builder.allAssigneeString(pull)
		mention = "=> " + name
	case 1:
		reaction = ":+1:"
		color = "warning"
		name := builder.reviewerString(pull, thumbsUppers)
		mention = "=> " + name
	default:
		reaction = ":+1::+1:"
		color = "good"
		name := "@" + builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
		mention = name + " *マージお願いします*"
	}

	var message = title + " " + reaction + "\n" + mention

	var attachment slackposter.Attachment
	attachment = slackposter.Attachment{
		Fallback: message,
		Text:     message,
		Color:    color,
		Fields:   []slackposter.Field{},
		MrkdwnIn: []string{"text", "fallback"},
	}

	return attachment
}

func NewMessageBuilder(gh GitHubAPI, usersManager UsersManager) MessageBuilder {
	return MessageBuilder{
		GitHubOwner:  gh.Owner,
		GitHubRepo:   gh.Repo,
		UsersManager: usersManager,
	}
}
