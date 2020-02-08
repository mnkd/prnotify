# prnotify
Notify GitHub pull requests to Slack incoming webhook.

## Installation
```
$ go get -u github.com/mnkd/prnotify
```

## How to make
```
$ git clone https://github.com/github/prnotify.git
$ cd prnotify
$ make
```

## Usage
```
$ prnotify
```

### Dry run
* A dry run will not send any message to Slack.

```
$ prnotify -d
```

## Configuration

### config.json
* Required
* Path `$HOME/.config/prnotify/config.json`

```json
{
  "app": {
    "use_holiday_jp": true,
    "use_dayoff": true
  },
  "github": {
    "access_token": "your access_token",
    "owner": "owner name",
    "repo": "repo name",
    "minimum_approved": 2
  },
  "slack_webhooks": {
    "team": "your-team",
    "channel": "#your-project",
    "username": "GitHub | Pull Requests",
    "icon_emoji": ":octocat:",
    "webhook_url": "https://hooks.slack.com/services/xxxxx/xxxxx/xxxxxx"
  }
}
```

### users.json
* Optional
* Path `$HOME/.config/prnotify/users.json`
* Dictionary for converting github username to slack username.

```json
{
  "github_username": "slack_username",
  "JohnnyAppleseed": "johnny",
  "TaroYamada": "taro"
}
```
