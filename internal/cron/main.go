// Package cron wraps github.com/robfig/cron/v3 to setup
// cronops tasks based on various plugins.
package cron

import (
	"fmt"

	"github.com/saimanwong/go-cronops/internal/config"
	"github.com/saimanwong/go-cronops/internal/logger"
	"github.com/saimanwong/go-cronops/internal/plugin"

	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron/v3"
)

var log = logger.New()

type Scheduler interface {
	Update(map[config.TaskID]config.Task)
	Status()

	stop()
	add(string, config.Task)
	del(cron.EntryID)
}

type Cron struct {
	sched *cron.Cron
	on    map[config.TaskID]cronTask
	off   map[config.TaskID]cronTask
}

type cronTask struct {
	id   cron.EntryID
	task config.Task
}

// Creates a new instance of Cron
func New() *Cron {
	return &Cron{
		sched: cron.New(),
		on:    map[config.TaskID]cronTask{},
		off:   map[config.TaskID]cronTask{},
	}
}

// Updates scheduler with new tasks
func (c *Cron) Update(tasks map[config.TaskID]config.Task) {
	// Reset scheduler
	c.reset()

	// Add OK tasks to enabled
	for id, t := range tasks {
		err := c.add(id, t)
		if err != nil {
			log.Error(err)
		}
	}

	// Start scheduler
	c.start()
}

func (c *Cron) Status() {
	for taskID, ct := range c.on {
		log.Infof("enabled %s (%s), next run %s", taskID, ct.task.Cron, c.sched.Entry(ct.id).Next)
	}
}

// Stops cron scheduler.
func (c *Cron) stop() {
	c.sched.Stop()
}

// Starts cron scheduler.
func (c *Cron) start() {
	c.sched.Start()
}

// Resets cron scheduler
func (c *Cron) reset() {
	for _, e := range c.sched.Entries() {
		c.sched.Remove(e.ID)
	}
	c.on = map[config.TaskID]cronTask{}
	c.off = map[config.TaskID]cronTask{}
}

// Adds task to scheduler
func (c *Cron) add(id config.TaskID, t config.Task) error {
	if t.Enabled {
		var job cron.Job

		// Add new plugins below
		switch t.Plugin {
		case "example":
			d := plugin.NewExample()
			err := mapstructure.Decode(t.Params, &d)
			if err != nil {
				return err
			}
			job = d
		default:
			return fmt.Errorf("no matching plugin %s for %s", t.Plugin, id)
		}

		// Only valid jobs will be enabled
		if job != nil {
			entryID, err := c.sched.AddJob(t.Cron, job)
			if err != nil {
				return fmt.Errorf("%s, probably bad cron expression", err)
			}

			c.on[id] = cronTask{
				id:   entryID,
				task: t,
			}
		}

		return nil
	}

	c.off[id] = cronTask{
		id:   -1,
		task: t,
	}
	return nil
}
