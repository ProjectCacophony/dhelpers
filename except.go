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
	if errD, ok := err.(*discordgo.RESTError); ok && errD.Message != nil {
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
			if msg != nil && cache.GetDiscord() != nil {
				// send message to discord, or add reaction if no message permission
				channelPermissions, chErr := cache.GetDiscord().UserChannelPermissions(event.BotUserID, msg.ChannelID)
				if chErr == nil {
					if channelPermissions&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages {
						errorMessage := err.Error()
						if errD, ok := err.(*discordgo.RESTError); ok && errD.Message != nil {
							errorMessage = errD.Message.Message
						}

						message := "**Something went wrong.** <a:ablobsadcloud:437572939701944322>\n```\nError: " + errorMessage + "\n```"
						if !dontLog {
							message += "I sent our top people to fix the issue as soon as possible. <a:ablobreach:437572330026434560>"
						}
						SendMessage( // nolint: errcheck
							msg.ChannelID,
							message,
						)
					} else if channelPermissions&discordgo.PermissionAddReactions == discordgo.PermissionAddReactions {
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
						cache.GetDiscord().MessageReactionAdd(msg.ChannelID, msg.ID, reactions[rand.Intn(len(reactions))]) // nolint: errcheck
					}
				}
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

// CheckErr panics if err is not nil
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
