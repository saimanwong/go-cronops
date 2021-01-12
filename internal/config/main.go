// Package config manages cronops overall config
package config

import (
	"fmt"

	"github.com/saimanwong/go-cronops/internal/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var log = logger.New()

type Task struct {
	Name      string
	Enabled   bool
	Cron      string
	Plugin    string
	Params    map[string]interface{}
	LogPrefix string
}

type ConfigChecker interface {
	reload()
	Watch(c chan bool)
}

type TaskID string
type Config struct {
	Tasks map[TaskID]Task
}

func New() *Config {
	cfg := &Config{}
	cfgPaths := []string{
		"/etc/cronops",
		"$HOME/.cronops",
		".",
		"./configs",
	}
	for _, p := range cfgPaths {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.reload()
	if err != nil {
		log.Error(err)
	}

	return cfg
}

// WatchConfig pushes true to channel if config is valid, otherwise false.
func (cfg *Config) Watch(c chan bool) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("config file changed: ", e.Name)
		err := cfg.reload()
		if err != nil {
			log.Error(err)
			c <- false
			return
		}
		c <- true
	})
}

// Reloads config if it is OK, otherwise returns error.
func (cfg *Config) reload() error {
	tasks := map[TaskID]Task{}

	// Parse tasks key
	// Note: keys become lower case
	// https://github.com/spf13/viper/pull/860
	err := viper.UnmarshalKey("tasks", &tasks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tasks, %w", err)
	}

	ok := validateTasks(tasks)
	if !ok {
		return fmt.Errorf("will not reload config utill it is fixed")
	}

	// Set default logprefix to TaskID
	for k, t := range tasks {
		t.Params["logPrefix"] = k
	}

	cfg.Tasks = tasks
	log.Infof("reloaded: %s", viper.ConfigFileUsed())
	return nil
}

// A function to validate that the config is OK.
func validateTasks(tasks map[TaskID]Task) bool {
	n := len(tasks)
	if n < 1 {
		log.Warn("no tasks specified")
	}

	log.Infof("found %d task(s)", n)
	ok := true
	for id, t := range tasks {
		if len(t.Name) < 1 {
			log.Warnf("task %s has no name", id)
			ok = false
		}
		if !t.Enabled {
			log.Warnf("task %s is disabled", id)
		}
		if len(t.Cron) < 1 {
			log.Warnf("task %s has no cron expression", id)
			ok = false
		}
		if len(t.Plugin) < 1 {
			log.Warnf("task %s has no plugin specified", id)
			ok = false
		}
		if len(t.Params) < 1 {
			log.Warnf("task %s has no param specified", id)
			ok = false
		}
	}

	return ok
}
