package dhelpers

import (
	"strings"

	"text/template"

	"github.com/dustin/go-humanize"
	"github.com/globalsign/mgo/bson"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/project-d-collab/dhelpers/cache"
	"gitlab.com/project-d-collab/dhelpers/mdb"
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
		// HumanizeComma adds commas to large numbers
		// example: {{HumanizeComma 1000}} => 1,000
		"NumberCommas": func(number int) string {
			return humanize.Comma(int64(number))
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
