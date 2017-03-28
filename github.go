package main

type GitHubAPI struct {
	AccessToken string
	Owner       string
	Repo        string
	Comment     struct {
		PerPage int
	}
}

func NewGitHubAPI(config Config) GitHubAPI {
	var gh = GitHubAPI{}
	gh.AccessToken = config.GitHub.AccessToken
	gh.Owner = config.GitHub.Owner
	gh.Repo = config.GitHub.Repo
	gh.Comment.PerPage = config.GitHub.Comment.PerPage
	return gh
}

func (gh GitHubAPI) BaseURL() string {
	return "https://api.github.com/repos/" + gh.Owner + "/" + gh.Repo
}
