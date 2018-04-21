package components

import (
	"os"

	"github.com/Seklfreak/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
	"gitlab.com/project-d-collab/dhelpers/cache"
	"gitlab.com/project-d-collab/discordrus"
)

// InitLogger initialises and caches the logging server
// will send warning and above messages to Discord if the DISCORD_LOGGING_WEBHOOK_URL environemnt variable is set
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

	// send warnings and above to discord if DISCORD_LOGGING_WEBHOOK_URL is set
	if os.Getenv("DISCORD_LOGGING_WEBHOOK_URL") != "" {
		log.AddHook(discordrus.NewHook(
			os.Getenv("DISCORD_LOGGING_WEBHOOK_URL"),
			logrus.WarnLevel,
			&discordrus.Opts{
				Username:         "Logging",
				Author:           "at " + service,
				DisableTimestamp: true,
			},
		))
	}

	cache.SetLogger(log.WithFields(logrus.Fields{"service": service}))
}
