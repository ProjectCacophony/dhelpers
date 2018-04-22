package dhelpers

import (
	"runtime"

	"fmt"

	"math/rand"
	"time"

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

	var dontLog bool
	// don't log discord permission errors
	if errD, ok := err.(*discordgo.RESTError); ok {
		if errD.Message.Code == discordgo.ErrCodeMissingPermissions ||
			errD.Message.Code == discordgo.ErrCodeMissingAccess ||
			errD.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			dontLog = true
		}
	}

	for _, errorHandlerType := range errorHandlers {
		switch errorHandlerType {
		case SentryErrorHandler:
			if dontLog {
				continue
			}

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
				// TODO: check channel permissions, if no permission to send message add reaction using addPermissionsErrorReaction
				SendMessage( // nolint: errcheck
					msg.ChannelID,
					"**Something went wrong.** <a:ablobsadcloud:437572939701944322>\n```\nError: "+err.Error()+"\n```I sent our top people to fix the issue as soon as possible. <a:ablobreach:437572330026434560>",
				)
			}
		}
	}

	if !dontLog {
		// log stacktrace
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)

		cache.GetLogger().Errorln(err.Error() + "\n\n" + string(buf[0:stackSize]))
	}
}

func addPermissionsErrorReaction(channelID, messageID string) {
	reactions := []string{
		":blobstop:317034621953114112",
		"a:ablobweary:394026914479865856",
		":googlespeaknoevil:317036753074651139",
		":notlikeblob:349342777978519562",
		"a:ablobcry:393869333740126219",
		"a:ablobfrown:394026913292615701",
		"a:ablobunamused:393869335573037057",
	}
	rand.Seed(time.Now().Unix())
	cache.GetDiscord().MessageReactionAdd(channelID, messageID, reactions[rand.Intn(len(reactions))]) // nolint: errcheck
}
