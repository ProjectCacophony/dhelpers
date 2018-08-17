package dhelpers

import (
	"runtime"

	"fmt"

	"math/rand"
	"time"

	"net"
	"syscall"

	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"gitlab.com/Cacophony/dhelpers/bucket"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/emoji"
	"gitlab.com/Cacophony/dhelpers/state"
)

// ErrorHandlerType is the Error Handler Type used for the EventContainer
type ErrorHandlerType string

// defines ErrorHandler types
const (
	SentryErrorHandler  ErrorHandlerType = "sentry"
	DiscordErrorHandler                  = "discord"
)

// defines ratelimiters for the ErrorHandlers
var (
	sentryLimiter  = bucket.NewBucket(5)
	discordLimiter = bucket.NewKeyBucket(3)
)

// HandleErrWith handles an error with the given error handles
// event can be nil
// currently supported ErrorHandlerTypes: SentryErrorHandler, and DiscordErrorHandler
func HandleErrWith(service string, err error, event *EventContainer, errorHandlers ...ErrorHandlerType) {
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

			if raven.ProjectID() == "" {
				continue
			}

			if !sentryLimiter.Allow() {
				continue
			}
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

		case DiscordErrorHandler:
			if event == nil || msg == nil || cache.GetEDiscord(event.BotUserID) == nil {
				continue
			}

			if !discordLimiter.Allow(msg.ChannelID) {
				continue
			}

			// send message to discord, or add reaction if no message permission
			channelPermissions, chErr := cache.GetEDiscord(event.BotUserID).UserChannelPermissions(event.BotUserID, msg.ChannelID)
			if chErr == nil {
				continue
			}

			if channelPermissions&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages {
				// send message if possible

				errorMessage := err.Error()
				if errD, ok := err.(*discordgo.RESTError); ok && errD.Message != nil {
					errorMessage = errD.Message.Message
				}

				message := "**Something went wrong.** " + emoji.GetWithout("sad") + "\n```\nError: " + errorMessage + "\n```"
				if !dontLog {
					message += "I sent our top people to fix the issue as soon as possible. " + emoji.GetWithout("reach")
				}
				event.SendMessage( // nolint: errcheck
					msg.ChannelID,
					message,
				)
			} else if channelPermissions&discordgo.PermissionAddReactions == discordgo.PermissionAddReactions {
				// try falling back to reaction if not possible

				reactions := []string{
					emoji.GetWithout("stop"),
					emoji.GetWithout("weary"),
					emoji.GetWithout("speaknoevil"),
					emoji.GetWithout("notlikethis"),
					emoji.GetWithout("cry"),
					emoji.GetWithout("frown"),
					emoji.GetWithout("unamused"),
				}
				rand.Seed(time.Now().Unix())
				cache.GetEDiscord(event.BotUserID).MessageReactionAdd(msg.ChannelID, msg.ID, reactions[rand.Intn(len(reactions))]) // nolint: errcheck
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

// HandleJobErrorWith handles a Job error, if errorHandlers is nil it will be sent to sentry
// currently supported ErrorHandlerTypes: SentryErrorHandler
func HandleJobErrorWith(service, job string, err error, errorHandlers ...ErrorHandlerType) {
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

// HandleHTTPErrorWith handles a HTTP error, if errorHandlers is nil it will be sent to sentry
// currently supported ErrorHandlerTypes: SentryErrorHandler
func HandleHTTPErrorWith(service string, request *http.Request, err error, errorHandlers ...ErrorHandlerType) {
	if errorHandlers == nil {
		errorHandlers = []ErrorHandlerType{SentryErrorHandler}
	}

	for _, errorHandlerType := range errorHandlers {
		switch errorHandlerType {
		case SentryErrorHandler:
			if raven.ProjectID() != "" {
				// send error to sentry
				raven.Capture(
					raven.NewPacket(
						err.Error(),
						raven.NewException(
							err,
							raven.NewStacktrace(4, 3, raven.IncludePaths()),
						),
						raven.NewHttp(request),
					),
					map[string]string{"service": service},
				)
			}
		}
	}

	if cache.GetLogger() != nil {
		// log stacktrace
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, false)

		cache.GetLogger().WithField("module", "http").Errorln(err.Error() + "\n\n" + string(buf[4:stackSize]))
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
