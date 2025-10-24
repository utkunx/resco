package services

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	translations map[string]string
	fallbackTranslations map[string]map[string]string
	translationMutex sync.RWMutex
	translationsLoaded bool
	fallbackLoaded bool
)

// LoadTranslations loads the Turkish to Chinese translations from the JSON file
func LoadTranslations() error {
	translationMutex.Lock()
	defer translationMutex.Unlock()

	if translationsLoaded {
		return nil
	}

	data, err := os.ReadFile("translate/tr-to-cn.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &translations)
	if err != nil {
		return err
	}

	translationsLoaded = true
	return nil
}

// LoadFallbackTranslations loads the fallback translations based on first 4 digits
func LoadFallbackTranslations() error {
	translationMutex.Lock()
	defer translationMutex.Unlock()

	if fallbackLoaded {
		return nil
	}

	data, err := os.ReadFile("translate/fallback-tr-to-cn.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &fallbackTranslations)
	if err != nil {
		return err
	}

	fallbackLoaded = true
	return nil
}

// Translate returns the Chinese translation for a Turkish text
// If no translation is found, it returns the original text
func Translate(turkishText string) string {
	translationMutex.RLock()
	defer translationMutex.RUnlock()

	if chineseText, exists := translations[turkishText]; exists {
		return chineseText
	}
	return turkishText
}

// TranslateWithFallback returns the Chinese translation using fallback logic
// First tries direct translation, then uses first 4 digits of itemCode for fallback
func TranslateWithFallback(turkishText string, itemCode string) string {
	translationMutex.RLock()
	defer translationMutex.RUnlock()

	// Try direct translation first
	if chineseText, exists := translations[turkishText]; exists {
		return chineseText
	}

	// If not found and itemCode is provided, try fallback based on first 4 digits
	if itemCode != "" && len(itemCode) >= 4 {
		prefix := itemCode[:4]
		if prefixTranslations, exists := fallbackTranslations[prefix]; exists {
			if chineseText, exists := prefixTranslations[turkishText]; exists {
				return chineseText
			}
		}
	}

	// Return original text if no translation found
	return turkishText
}

// ApplyTranslationsToBOM applies Chinese translations to BOM results
func ApplyTranslationsToBOM(results []BOMResult) []BOMResult {
	translatedResults := make([]BOMResult, len(results))

	for i, result := range results {
		translatedResults[i] = result

		// Translate parent name with fallback using parent number
		translatedResults[i].AD = TranslateWithFallback(result.AD, result.BOMRecCode)

		// Translate child name if it exists, with fallback using child number
		if result.SubItemName != nil && *result.SubItemName != "" {
			translated := TranslateWithFallback(*result.SubItemName, result.BOMRecKaynakCode)
			translatedResults[i].SubItemName = &translated
		}
	}

	return translatedResults
}