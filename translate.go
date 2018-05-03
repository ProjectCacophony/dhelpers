package dhelpers

import (
	"strings"

	"text/template"

	"time"

	"github.com/dustin/go-humanize"
	"github.com/globalsign/mgo/bson"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/mdb"
)

var (
	// the additional functions to use in the template engine
	translationFuncs = template.FuncMap{
		// ObjectID returns a humanised version of an object ID
		// example: {{ObjectID id}} => 5ae3a59624f8753dba273792
		"ObjectID": func(id bson.ObjectId) string {
			return mdb.IDToHuman(id)
		},
		// MarkdownLinkEscape returns a given link escaped to be used in Markdown
		// example: {{MarkdownLinkEscape "https://example.org/A+(B)"}} => https://example.org/A+%28B%29
		"MarkdownLinkEscape": func(text string) string {
			return EscapeLinkForMarkdown(text)
		},
		// NumberCommas adds commas to large numbers
		// example: {{HumanizeComma 1000}} => 1,000
		"NumberCommas": func(number int) string {
			return humanize.Comma(int64(number))
		},
		// HumanizeTime formats time human readable
		// example: {{HumanizeTime time.Sub(10*time.Minute)}} => 10 minutes ago
		"HumanizeTime": func(theTime time.Time) string {
			return humanize.Time(theTime)
		},
		// Prefix returns the prefix for a GuildID
		// example: {{Prefix 339227598544568340}} => /
		"Prefix": func(guildID string) string {
			return GetPrefix(guildID)
		},
		// PrefixE returns the prefix for an EventContainer
		// currently only works with MessageCreate, MessageUpdate, and MessageDelete events
		// example: {{PrefixE event}} => /
		"PrefixE": func(event EventContainer) string {
			switch event.Type {
			case MessageCreateEventType:
				return GetPrefix(event.MessageCreate.GuildID)
			case MessageUpdateEventType:
				return GetPrefix(event.MessageUpdate.GuildID)
			case MessageDeleteEventType:
				return GetPrefix(event.MessageDelete.GuildID)
			}
			return defaultPrefix
		},
	}
)

// T returns the translation for the given message ID
// Example: T("HelloWorld")
func T(messageID string) (result string) {
	if cache.GetLocalizationBundle() == nil {
		return messageID
	}

	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
			result = messageID
		}
	}()

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
		Funcs:     translationFuncs,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") { // ignore message not found errors
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
		}
		return messageID
	}
	return translation
}

// Tf returns the translation for the given message ID applying the fields
// Example: Tf("HelloWorld", "key", "value")
func Tf(messageID string, fields ...interface{}) (result string) {
	if cache.GetLocalizationBundle() == nil {
		return messageID
	}

	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
			result = messageID
		}
	}()

	// create map out of fields
	data := make(map[interface{}]interface{})
	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			data[fields[i]] = fields[i+1]
		}
	}

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		Funcs:        translationFuncs,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") { // ignore message not found errors
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
		}
		return messageID
	}
	return translation
}

// Tfc returns the translation for the given message ID applying the fields and pluralization count
// Example: Tfc("HelloWorld", 3, "key", "value")
func Tfc(messageID string, count int, fields ...interface{}) (result string) {
	if cache.GetLocalizationBundle() == nil {
		return messageID
	}

	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
			result = messageID
		}
	}()

	// create map out of fields
	data := make(map[interface{}]interface{})
	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			data[fields[i]] = fields[i+1]
		}
	}

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		PluralCount:  count,
		Funcs:        translationFuncs,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") { // ignore message not found errors
			cache.GetLogger().WithField("module", "translate").Errorln(err.(error).Error())
		}
		return messageID
	}
	return translation
}

// T returns the translation for the given message ID, the event variable is being set
// Example: T("HelloWorld")
func (event EventContainer) T(messageID string) (result string) {
	return Tf(messageID, "event", event)
}

// Tf returns the translation for the given message ID applying the fields, the event variable is being set
// Example: Tf("HelloWorld", "key", "value")
func (event EventContainer) Tf(messageID string, fields ...interface{}) (result string) {
	return Tf(messageID, append(fields, "event", event)...)
}

// Tfc returns the translation for the given message ID applying the fields and pluralization count, the event variable is being set
// Example: Tfc("HelloWorld", 3, "key", "value")
func (event EventContainer) Tfc(messageID string, count int, fields ...interface{}) (result string) {
	return Tfc(messageID, count, append(fields, "event", event)...)
}
