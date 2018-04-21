package dhelpers

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// T returns the translation for the given message ID
// Example: T("HelloWorld")
func T(messageID string) (result string) {
	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			result = messageID
		}
	}()

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID:   messageID,
		PluralCount: 1,
	})
	if err != nil {
		return messageID
	}
	return translation
}

// Tf returns the translation for the given message ID applying the fields
// Example: Tf("HelloWorld", "key", "value")
func Tf(messageID string, fields ...string) (result string) {
	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			result = messageID
		}
	}()

	// create map out of fields
	data := make(map[string]string)
	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			data[fields[i]] = fields[i+1]
		}
	}

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		PluralCount:  1,
	})
	if err != nil {
		return messageID
	}
	return translation
}

// Tfc returns the translation for the given message ID applying the fields and pluralization count
// Example: Tfc("HelloWorld", 3, "key", "value")
func Tfc(messageID string, count int, fields ...string) (result string) {
	// on panic return message ID
	defer func() {
		err := recover()
		if err != nil {
			result = messageID
		}
	}()

	// create map out of fields
	data := make(map[string]string)
	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			data[fields[i]] = fields[i+1]
		}
	}

	translation, err := i18n.NewLocalizer(cache.GetLocalizationBundle(), "en").Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		PluralCount:  count,
	})
	if err != nil {
		return messageID
	}
	return translation
}
