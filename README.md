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
- A dry run will not send any message to Slack.

```
$ prnotify -d
```

## Configuration

### config.json
- Required
- Path `$HOME/.config/prnotify/config.json`

```
{
  "app": {
    "use_holiday_jp": true,
    "use_dayoff": true
  },
  "github": {
    "access_token": "your access_token",
    "owner": "owner name",
    "repo": "repo name",
    "comment": {
       "per_page": 100
    }
  },
  "slack_webhooks": [
    {
      "team": "your-team",
      "channel": "#your-project",
      "username": "GitHub | Pull Requests",
      "icon_emoji": ":octocat:",
      "webhook_url": "https://hooks.slack.com/services/xxxxx/xxxxx/xxxxxx"
    }
  ]
}
```

### users.json
- Optional
- Path `$HOME/.config/prnotify/users.json`
- Dictionary for converting github username to slack username.

```
{
  "github_username": "slack_username",
  "JohnnyAppleseed": "johnny",
  "TaroYamada": "taro"
}
```
