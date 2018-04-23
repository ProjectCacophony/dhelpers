package components

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/project-d-collab/dhelpers/cache"
	"golang.org/x/text/language"
)

// InitTranslator initialises and caches a translation bundle
func InitTranslator(files []string) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, file := range files {
		cache.GetLogger().Infoln("Loaded " + file)
		bundle.MustLoadMessageFile(file)
	}

	cache.SetLocalizationBundle(bundle)
}
