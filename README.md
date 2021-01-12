# cronops

cronops is an (monolith, for now) app that tries to automate away manual processes with modular plugins and periodic tasks.

## Getting started

### Installing

```
$ go get -u github.com/saimanwong/go-cronops/cmd/cronops
```

## Built with

### Cronops

* [spf13/viper](https://github.com/spf13/viper) - Used to manage Go configuration dynamically
* [fsnotify/fsnotify](github.com/fsnotify/fsnotify) - Watches configuration file and notifies cronops
* [robfig/cron](https://github.com/robfig/cron) - Schedule cron jobs in Go
* [sirupsen/logrus](github.com/sirupsen/logrus) - For logging
* [mitchellh/mapstructure](github.com/mitchellh/mapstructure) - Decode maps in Go

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
