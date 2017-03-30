package main

import (
	"fmt"
	"github.com/mnkd/slackposter"
	// "strings"
)

type MessageBuilderForReviews struct {
	GitHubOwner     string
	GitHubRepo      string
	UsersManager    UsersManager
	MinimumApproved int
}

func (builder MessageBuilderForReviews) fieldTitleString(pull PullRequest) string {
	return fmt.Sprintf("#%d", pull.Number)
}

func (builder MessageBuilderForReviews) fieldValueString(pull PullRequest) string {
	return fmt.Sprintf("<%s|%s> by %s", pull.HTMLURL, pull.Title, pull.User.Login)
}

func (builder MessageBuilderForReviews) allReviewersString(pull PullRequest, reviewers []RequestedReviewer) string {
	if len(reviewers) == 0 {
		name := builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
		return "@" + name + " *Reviewers の指定をお願いします*"
	}

	var str = ""
	for _, reviewer := range reviewers {
		name := builder.UsersManager.ConvertGitHubToSlack(reviewer.Login)
		str += "@" + name + " "
	}
	return str
}

func (builder MessageBuilderForReviews) reviewerString(pull PullRequest, reviewers []string) string {
	var str = ""
	for _, reviewer := range reviewers {
		login := builder.UsersManager.ConvertGitHubToSlack(reviewer)
		str += "@" + login + " "
	}
	return str
}

func (builder MessageBuilderForReviews) BudildSummary(pullsCount int) string {
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

type Usernames []string

func (usernames Usernames) isContain(username string) bool {
	for _, u := range usernames {
		if u == username {
			return true
		}
	}
	return false
}

func (builder MessageBuilderForReviews) BuildField(pull PullRequest, reviewers []RequestedReviewer, reviews []Review) (slackposter.Field, AttachmentType) {
	var approvedUsers Usernames
	var changeRequestedUsers Usernames

	for _, review := range reviews {
		username := builder.UsersManager.ConvertGitHubToSlack(review.User.Login)
		if review.IsApproved() {
			approvedUsers = append(approvedUsers, username)

		} else if review.IsRequestedChanged() {
			changeRequestedUsers = append(changeRequestedUsers, username)
		}
	}

	var unreviewUsers Usernames
	for _, reviewer := range reviewers {
		username := builder.UsersManager.ConvertGitHubToSlack(reviewer.Login)
		if !approvedUsers.isContain(username) && !changeRequestedUsers.isContain(username) {
			unreviewUsers = append(unreviewUsers, username)
		}
	}

	fmt.Println("approvedUsers:", approvedUsers)
	fmt.Println("changeRequestedUsers:", changeRequestedUsers)
	fmt.Println("unreviewUsers:", unreviewUsers)

	var attachmentType AttachmentType
	title := builder.fieldTitleString(pull)
	value := builder.fieldValueString(pull)
	pullUsername := builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
	name := ""

	if len(reviewers) == 0 {
		attachmentType = REVIEWERS
		name = "@" + pullUsername
	} else {
		if len(approvedUsers) >= builder.MinimumApproved {
			attachmentType = MERGE
			name = "@" + pullUsername
		} else if len(changeRequestedUsers) > 0 {
			attachmentType = REQUEST_CHANGE
			name = "@" + pullUsername
		} else {
			attachmentType = REVIEW_ONE
			name = builder.reviewerString(pull, unreviewUsers)
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

func (builder MessageBuilderForReviews) BuildAttachmentWithType(attachmentType AttachmentType) slackposter.Attachment {
	var color, message string
	switch attachmentType {
	case MERGE:
		color = "good"
		message = ":+1::+1: *マージお願いします*"
	case REVIEW_ONE:
		color = "warning"
		message = ":smiley: *レビューお願いします！*"
	case REVIEWERS:
		color = "danger"
		message = ":sweat_smile: *Reviewers の指定をお願いします！*"
	case REQUEST_CHANGE:
		color = "danger"
		message = ":wink: *コードの修正をお願いします！*"
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

func NewMessageBuilderForReviews(gh GitHubAPI, usersManager UsersManager, config Config) MessageBuilderForReviews {
	return MessageBuilderForReviews{
		GitHubOwner:     gh.Owner,
		GitHubRepo:      gh.Repo,
		UsersManager:    usersManager,
		MinimumApproved: config.GitHub.MinimumApproved,
	}
}
