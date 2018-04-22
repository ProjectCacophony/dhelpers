package dhelpers

import (
	"runtime"

	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// ErrorHandlerType is the Error Handler Type used for the EventContainer
type ErrorHandlerType string

// defines ErrorHandler types
const (
	SentryErrorHandler  ErrorHandlerType = "sentry"
	DiscordErrorHandler                  = "discord"
)

// HandleErrWith handles an error with the given error handles
// event can be nil
func HandleErrWith(service string, err error, errorHandlers []ErrorHandlerType, event *EventContainer) {
	var msg *discordgo.Message
	if event.MessageCreate != nil {
		msg = event.MessageCreate.Message
	}
	if event.MessageUpdate != nil {
		msg = event.MessageUpdate.Message
	}

	for _, errorHandlerType := range errorHandlers {
		switch errorHandlerType {
		case SentryErrorHandler:
			if raven.ProjectID() != "" {
				// send error to sentry
				data := map[string]string{"service": service}
				if msg != nil {
					data["MessageID"] = msg.ID
					data["AuthorID"] = msg.Author.ID
					data["ChannelID"] = msg.ChannelID
					data["Content"] = msg.Content
					data["Timestamp"] = string(msg.Timestamp)
				}

				raven.CaptureError(fmt.Errorf(spew.Sdump(err)), data)
			}
		case DiscordErrorHandler:
			if msg != nil {
				// send message to discord
				SendMessage( // nolint: errcheck
					msg.ChannelID,
					"Something went wrong. <a:ablobfrown:394026913292615701>\n```\n"+err.Error()+"\n```",
				)
			}
		}
	}

	// log stacktrace
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, false)

	cache.GetLogger().Errorln(string(buf[0:stackSize]))
}
