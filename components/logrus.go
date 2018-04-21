package components

import (
	"os"

	"github.com/Sirupsen/logrus"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

func InitLogger() {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "02-01-06 15:04:05.000"}
	log.Hooks = make(logrus.LevelHooks)

	cache.SetLogger(log)
}
