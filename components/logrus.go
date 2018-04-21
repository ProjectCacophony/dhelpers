package components

import (
	"os"

	"github.com/Seklfreak/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

func InitLogger(service string) {
	format := new(prefixed.TextFormatter)
	format.TimestampFormat = "02-01-06 15:04:05.000"
	format.FullTimestamp = true
	format.ForceColors = true
	format.SpacePadding = 2

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	log.Formatter = format

	// log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "02-01-06 15:04:05.000"}
	log.Hooks = make(logrus.LevelHooks)

	cache.SetLogger(log.WithFields(logrus.Fields{"service": service}))
}
