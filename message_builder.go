package main

import (
	"fmt"

	"github.com/mnkd/slackposter"
	// "strings"
)

type Usernames []string

func (usernames Usernames) isContain(username string) bool {
	for _, u := range usernames {
		if u == username {
			return true
		}
	}
	return false
}

func UsernameFromRequestedReviewers(requestedReviewers []PullRequestUser) []string {
	var array []string
	for _, r := range requestedReviewers {
		array = append(array, r.Login)
	}
	return array
}

type MessageBuilder struct {
	GitHubOwner     string
	GitHubRepo      string
	UsersManager    UsersManager
	MinimumApproved int
}

func (builder MessageBuilder) fieldTitleString(pull PullRequest) string {
	return fmt.Sprintf("#%d", pull.Number)
}

func (builder MessageBuilder) fieldValueString(pull PullRequest) string {
	return fmt.Sprintf("<%s|%s> by %s", pull.HTMLURL, pull.Title, pull.User.Login)
}

func (builder MessageBuilder) allReviewersString(pull PullRequest, requestedReviewers []PullRequestUser) string {
	if len(requestedReviewers) == 0 {
		name := builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
		return "@" + name + " *Reviewers の指定をお願いします*"
	}

	var str = ""
	for _, reviewer := range requestedReviewers {
		name := builder.UsersManager.ConvertGitHubToSlack(reviewer.Login)
		str += "@" + name + " "
	}
	return str
}

func (builder MessageBuilder) reviewerString(pull PullRequest, reviewers []string) string {
	var str = ""
	for _, reviewer := range reviewers {
		login := builder.UsersManager.ConvertGitHubToSlack(reviewer)
		str += "@" + login + " "
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

func (builder MessageBuilder) isAssigneeReview(review Review, pull PullRequest) bool {
	for _, assignee := range pull.Assignees {
		if assignee.Login == review.User.Login {
			return true
		}
	}
	return false
}

func (builder MessageBuilder) isRequestedReviewerReview(review Review, pull PullRequest) bool {
	for _, requestedReviewer := range pull.RequestedReviewers {
		if requestedReviewer.Login == review.User.Login {
			return true
		}
	}
	return false
}

func (builder MessageBuilder) BuildField(pull PullRequest, reviews []Review) (slackposter.Field, AttachmentType) {
	var approvedUsers Usernames    // User who have approved this pull request.
	var notApprovedUsers Usernames // User who have not approved this pull request.
	requestedReviewers := UsernameFromRequestedReviewers(pull.RequestedReviewers)

	for _, review := range reviews {
		username := builder.UsersManager.ConvertGitHubToSlack(review.User.Login)

		if review.IsApproved() {
			approvedUsers = append(approvedUsers, username)
		} else if builder.isAssigneeReview(review, pull) == false &&
			builder.isRequestedReviewerReview(review, pull) == false &&
			notApprovedUsers.isContain(username) == false {
			notApprovedUsers = append(notApprovedUsers, username)
		}
	}

	fmt.Println("requestedReviewers:", requestedReviewers)
	fmt.Println("approvedUsers:", approvedUsers)
	fmt.Println("notApprovedUsers:", notApprovedUsers)

	var attachmentType AttachmentType
	title := builder.fieldTitleString(pull)
	value := builder.fieldValueString(pull)
	pullUsername := builder.UsersManager.ConvertGitHubToSlack(pull.User.Login)
	name := ""

	if len(approvedUsers) >= builder.MinimumApproved {
		attachmentType = MERGE
		name = "@" + pullUsername

	} else if len(requestedReviewers) > 0 {
		attachmentType = REVIEW
		name = builder.reviewerString(pull, requestedReviewers) + " " + builder.reviewerString(pull, notApprovedUsers)

	} else if len(requestedReviewers) == 0 && len(reviews) == 0 {
		attachmentType = ASSIGN_REVIEWER
		name = "@" + pullUsername

	} else {
		attachmentType = CHECK
		name = "@" + pullUsername + " " + builder.reviewerString(pull, notApprovedUsers)
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
	case REVIEW:
		color = "warning"
		message = ":smiley: *引き続きレビューお願いします！*"
	case CHECK:
		color = "warning"
		message = ":wink: *進捗／Reviewers の再指定／APPROVED し忘れ／など確認お願いします！*"
	case ASSIGN_REVIEWER:
		color = "danger"
		message = ":sweat_smile: *Reviewers の指定をお願いします！*"
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

func NewMessageBuilder(gh GitHubAPI, usersManager UsersManager, config Config) MessageBuilder {
	return MessageBuilder{
		GitHubOwner:     gh.Owner,
		GitHubRepo:      gh.Repo,
		UsersManager:    usersManager,
		MinimumApproved: config.GitHub.MinimumApproved,
	}
}
