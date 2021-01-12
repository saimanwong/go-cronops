# cronops

cronops is an (monolith, for now) app that tries to automate away manual processes with modular plugins and periodic tasks.

## Getting started

### Installing

```
$ go get -u github.com/saimanwong/go-cronops/cmd/cronops
```

### Example

```
$ cronops
[2021-01-11 22:15:00]  INFO found 3 task(s)
[2021-01-11 22:15:00]  WARN task example2 is disabled
[2021-01-11 22:15:00]  INFO reloaded: /Users/username/.cronops/config.yml
[2021-01-11 22:15:00]  INFO enabled example (@every 5s), next run 2021-01-11 22:15:05 +0100 CET
[2021-01-11 22:15:00]  INFO enabled jiratoslack (@every 5s), next run 2021-01-11 22:15:05 +0100 CET
[2021-01-11 22:15:05]  INFO example: hello world
[2021-01-11 22:15:05]  INFO jiratoslack: ingored issueKey TEST-6
[2021-01-11 22:15:05]  INFO jiratoslack: ignored summary ignore
[2021-01-11 22:15:07]  INFO jiratoslack: posted the first message to slack channel CXXXXXXXXXX
[2021-01-11 22:15:07]  INFO jiratoslack: posted issue TEST-8 to slack thread
[2021-01-11 22:15:07]  INFO jiratoslack: posted issue TEST-7 to slack thread
[2021-01-11 22:15:07]  INFO jiratoslack: posted issue TEST-3 to slack thread
[2021-01-11 23:02:22]  INFO config file changed: /Users/username/.cronops/config.yml
[2021-01-11 23:02:22]  INFO found 3 task(s)
[2021-01-11 23:02:22]  WARN task example2 is disabled
[2021-01-11 23:02:22]  WARN task jiratoslack is disabled
[2021-01-11 23:02:22]  INFO reloaded: /Users/username/.cronops/config.yml
[2021-01-11 23:02:22]  INFO enabled example (@every 5s), next run 2021-01-11 23:02:27 +0100 CET
[2021-01-11 23:02:27]  INFO example: hello world
```

## Built with

### Cronops

* [spf13/viper](https://github.com/spf13/viper) - Used to manage Go configuration dynamically
* [fsnotify/fsnotify](github.com/fsnotify/fsnotify) - Watches configuration file and notifies cronops
* [robfig/cron](https://github.com/robfig/cron) - Schedule cron jobs in Go
* [sirupsen/logrus](github.com/sirupsen/logrus) - For logging
* [mitchellh/mapstructure](github.com/mitchellh/mapstructure) - Decode maps in Go

### [JiraSlack Plugin](internal/plugin/jiraslack)

* [andygrunwald/go-jira](https://github.com/andygrunwald/go-jira) - Used to search for issues in Jira
* [slack-go/slack](https://github.com/slack-go/slack) - Slack library for Go

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
