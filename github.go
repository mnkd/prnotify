package main

type GitHubAPI struct {
	AccessToken string
	Owner       string
	Repo        string
}

func NewGitHubAPI(config Config) GitHubAPI {
	var gh = GitHubAPI{}
	gh.AccessToken = config.GitHub.AccessToken
	gh.Owner = config.GitHub.Owner
	gh.Repo = config.GitHub.Repo
	return gh
}

func (gh GitHubAPI) BaseURL() string {
	return "https://api.github.com/repos/" + gh.Owner + "/" + gh.Repo
}
