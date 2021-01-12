// Binary cronops.
package main

import (
	"github.com/saimanwong/go-cronops/internal/config"
	"github.com/saimanwong/go-cronops/internal/cron"
)

func main() {
	cfg := config.New()
	cfgChan := make(chan bool)
	cfg.Watch(cfgChan)
	updateTasks(cfgChan, cfg)
}

// updateTasks updates cron scheduler with new tasks.
func updateTasks(cfgChan chan bool, cfg *config.Config) {
	c := cron.New()
	c.Update(cfg.Tasks)
	c.Status()
	for {
		select {
		case ok := <-cfgChan:
			if ok {
				c.Update(cfg.Tasks)
				c.Status()
			}
		}
	}
}
