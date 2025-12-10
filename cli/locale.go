package cli

import (
	"os"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/bg"
	"github.com/go-playground/locales/cs"
	"github.com/go-playground/locales/da"
	"github.com/go-playground/locales/de"
	"github.com/go-playground/locales/el"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/et_EE"
	"github.com/go-playground/locales/fi"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/he"
	"github.com/go-playground/locales/hi"
	"github.com/go-playground/locales/hr"
	"github.com/go-playground/locales/hu"
	"github.com/go-playground/locales/id"
	"github.com/go-playground/locales/it"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/lt"
	"github.com/go-playground/locales/lv"
	"github.com/go-playground/locales/nb"
	"github.com/go-playground/locales/nl"
	"github.com/go-playground/locales/pl"
	"github.com/go-playground/locales/pt"
	"github.com/go-playground/locales/ro"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/sk"
	"github.com/go-playground/locales/sl"
	"github.com/go-playground/locales/sr"
	"github.com/go-playground/locales/sv"
	"github.com/go-playground/locales/th"
	"github.com/go-playground/locales/tr"
	"github.com/go-playground/locales/uk"
	"github.com/go-playground/locales/vi"
	"github.com/go-playground/locales/zh"
)

// getSystemLocale detects the system locale and returns the appropriate translator.
func getSystemLocale() locales.Translator {
	// Check common environment variables for locale
	locale := os.Getenv("LC_TIME")
	if locale == "" {
		locale = os.Getenv("LC_ALL")
	}
	if locale == "" {
		locale = os.Getenv("LANG")
	}

	// Extract language code (e.g., "en_US.UTF-8" -> "en")
	lang := strings.ToLower(locale)
	if idx := strings.Index(lang, "_"); idx > 0 {
		lang = lang[:idx]
	}

	// Return appropriate translator
	switch lang {
	case "ar": // Arabic
		return ar.New()
	case "bg": // Bulgarian
		return bg.New()
	case "cs": // Czech
		return cs.New()
	case "da": // Danish
		return da.New()
	case "de": // German
		return de.New()
	case "el": // Greek
		return el.New()
	case "es": // Spanish
		return es.New()
	case "et": // Estonian
		return et_EE.New()
	case "fi": // Finnish
		return fi.New()
	case "fr": // French
		return fr.New()
	case "he": // Hebrew
		return he.New()
	case "hi": // Hindi
		return hi.New()
	case "hr": // Croatian
		return hr.New()
	case "hu": // Hungarian
		return hu.New()
	case "id": // Indonesian
		return id.New()
	case "it": // Italian
		return it.New()
	case "ja": // Japanese
		return ja.New()
	case "ko": // Korean
		return ko.New()
	case "lt": // Lithuanian
		return lt.New()
	case "lv": // Latvian
		return lv.New()
	case "nb", "no": // Norwegian (Bokmål)
		return nb.New()
	case "nl": // Dutch
		return nl.New()
	case "pl": // Polish
		return pl.New()
	case "pt": // Portuguese
		return pt.New()
	case "ro": // Romanian
		return ro.New()
	case "ru": // Russian
		return ru.New()
	case "sk": // Slovak
		return sk.New()
	case "sl": // Slovenian
		return sl.New()
	case "sr": // Serbian
		return sr.New()
	case "sv": // Swedish
		return sv.New()
	case "th": // Thai
		return th.New()
	case "tr": // Turkish
		return tr.New()
	case "uk": // Ukrainian
		return uk.New()
	case "vi": // Vietnamese
		return vi.New()
	case "zh": // Chinese
		return zh.New()
	default: // English as fallback
		return en.New()
	}
}
