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

func (builder MessageBuilder) fieldTitleString(pull PullRequest) string {
	return fmt.Sprintf("#%d", pull.Number)
}

func (builder MessageBuilder) fieldValueString(pull PullRequest) string {
	return fmt.Sprintf("<%s|%s> by %s", pull.HTMLURL, pull.Title, pull.User.Login)
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

func (builder MessageBuilder) reviewerString(pull PullRequest, reviewdUsers []string) string {
	var str = ""
	for _, assignee := range pull.Assignees {
		assigneeLogin := builder.UsersManager.ConvertGitHubToSlack(assignee.Login)
		found := false
		for _, reviewdUser := range reviewdUsers {
			if assigneeLogin == reviewdUser {
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

func (builder MessageBuilder) BuildField(pull PullRequest, comments []Comment) (slackposter.Field, AttachmentType) {
	var reviewdUsers []string
	for _, comment := range comments {
		if strings.Contains(comment.Body, ":+1:") || strings.Contains(comment.Body, "👍") {
			username := builder.UsersManager.ConvertGitHubToSlack(comment.User.Login)
			reviewdUsers = append(reviewdUsers, username)
		}
	}

	var attachmentType AttachmentType
	title := builder.fieldTitleString(pull)
	value := builder.fieldValueString(pull)
	name := ""

	if len(pull.Assignees) == 0 {
		attachmentType = ASSIGNEE
		name = "@" + builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
	} else {
		switch len(reviewdUsers) {
		case 0:
			attachmentType = REVIEW_TWO
			name = builder.allAssigneeString(pull)
		case 1:
			attachmentType = REVIEW_ONE
			name = builder.reviewerString(pull, reviewdUsers)
		default:
			attachmentType = MERGE
			name = "@" + builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
		}
	}

	value = value + " => " + name

	field := slackposter.Field{
		Title: title,
		Value: value,
		Short: false,
	}

	return field, attachmentType
}

func (builder MessageBuilder) BuildAttachmentWithType(attachmentType AttachmentType) slackposter.Attachment {
	var color, message string
	switch attachmentType {
	case MERGE:
		color = "good"
		message = ":+1::+1: *マージお願いします*"
	case REVIEW_ONE:
		color = "warning"
		message = ":+1: *レビューお願いします！*"
	case REVIEW_TWO:
		color = "danger"
		message = ":wink: *レビューお願いします！*"
	case ASSIGNEE:
		color = "danger"
		message = ":sweat_smile: *Assignee の指定をお願いします！*"
	}

	var attachment slackposter.Attachment
	attachment = slackposter.Attachment{
		Fallback: message,
		Text:     message,
		Color:    color,
		Fields:   []slackposter.Field{},
		MrkdwnIn: []string{"text", "fallback"},
		// AuthorIcon: "https://assets-cdn.github.com/images/modules/logos_page/GitHub-Mark.png",
		// AuthorLink: "https://github.com/",
		// AuthorName: "GitHub",
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
