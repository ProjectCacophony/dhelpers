package dhelpers

import (
	"runtime"

	"fmt"

	"math/rand"
	"time"

	"net"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"gitlab.com/project-d-collab/dhelpers/cache"
	"gitlab.com/project-d-collab/dhelpers/state"
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
// currently supported ErrorHandlerTypes: SentryErrorHandler, and DiscordErrorHandler
func HandleErrWith(service string, err error, errorHandlers []ErrorHandlerType, event *EventContainer) {
	var msg *discordgo.Message
	if event != nil {
		if event.MessageCreate != nil {
			msg = event.MessageCreate.Message
		}
		if event.MessageUpdate != nil {
			msg = event.MessageUpdate.Message
		}
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
	// don't log state errors
	if err == state.ErrStateNotFound ||
		err == state.ErrTargetWrongServer ||
		err == state.ErrTargetWrongType {
		dontLog = true
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
			if event != nil && msg != nil && cache.GetEDiscord(event.BotUserID) != nil {
				// send message to discord, or add reaction if no message permission
				channelPermissions, chErr := cache.GetEDiscord(event.BotUserID).UserChannelPermissions(event.BotUserID, msg.ChannelID)
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
						event.SendMessage( // nolint: errcheck
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
						cache.GetEDiscord(event.BotUserID).MessageReactionAdd(msg.ChannelID, msg.ID, reactions[rand.Intn(len(reactions))]) // nolint: errcheck
					}
				}
			}
		}
	}

	if !dontLog && cache.GetLogger() != nil {
		// log stacktrace
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)

		cache.GetLogger().Errorln(err.Error() + "\n\n" + string(buf[0:stackSize]))
	}
}

// RecoverLog can be recovered to, all errors will be logged
func RecoverLog() {
	err := recover()
	if err != nil {
		// handle errors
		LogError(err.(error))
	}
}

// LogError sends an error to sentry and logs it, can be nil
func LogError(err error) {
	if err == nil {
		return
	}

	if cache.GetLogger() != nil {
		// log stacktrace
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)

		cache.GetLogger().Errorln(err.Error() + "\n\n" + string(buf[0:stackSize]))
	}

	if raven.ProjectID() != "" {
		// send error to sentry
		raven.CaptureError(fmt.Errorf(spew.Sdump(err)), nil)
	}
}

// HandleJobError handles a Job error, if errorHandlers is nil it will be sent to sentry
// currently supported ErrorHandlerTypes: SentryErrorHandler
func HandleJobError(service, job string, err error, errorHandlers []ErrorHandlerType) {
	if errorHandlers == nil {
		errorHandlers = []ErrorHandlerType{SentryErrorHandler}
	}

	for _, errorHandlerType := range errorHandlers {
		switch errorHandlerType {
		case SentryErrorHandler:
			if raven.ProjectID() != "" {
				// send error to sentry
				data := map[string]string{"service": service, "job": job}

				raven.CaptureError(fmt.Errorf(spew.Sdump(err)), data)
			}
		}
	}

	if cache.GetLogger() != nil {
		// log stacktrace
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)

		cache.GetLogger().WithField("job", job).Errorln(err.Error() + "\n\n" + string(buf[0:stackSize]))
	}
}

// CheckErr panics if err is not nil
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// IsNetworkErr returns true if an error is network related
func IsNetworkErr(err error) (is bool) {
	if err == nil {
		return false
	}

	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		return true
	}

	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			return true
		} else if t.Op == "read" {
			return true
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return true
		}
	}

	return false
}
