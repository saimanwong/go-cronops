package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func New() *logrus.Logger {
	return &logrus.Logger{
		Out: os.Stdout,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}
}

func NewLogEntry() *logrus.Entry {
	return &logrus.Entry{}
}
