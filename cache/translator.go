package cache

import (
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	localizationBundle      *i18n.Bundle
	localizationBundleMutex sync.RWMutex
)

// SetLocalizationBundle caches a localization bundle for future use
func SetLocalizationBundle(s *i18n.Bundle) {
	localizationBundleMutex.Lock()
	localizationBundle = s
	localizationBundleMutex.Unlock()
}

// GetLocalizationBundle returns a cached localization bundle
func GetLocalizationBundle() *i18n.Bundle {
	localizationBundleMutex.RLock()
	defer localizationBundleMutex.RUnlock()

	return localizationBundle
}
