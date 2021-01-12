package example

import (
	l "github.com/saimanwong/go-cronops/internal/logger"
)

var (
	logger = l.New()
	log    = l.NewLogEntry()
)

type Example struct {
	LogPrefix string
	Msg       string
}

func (e *Example) Run() {
	log = logger.WithField("prefix", e.LogPrefix)

	// Implementation after this
	log.Info(e.Msg)
}
