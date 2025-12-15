package i18n

import "strings"

type Translator struct {
	defaultLang string
	messages    map[string]map[string]string
}

func New(defaultLang string) *Translator {
	return &Translator{
		defaultLang: defaultLang,
		messages: map[string]map[string]string{
			"en": {
				"bad_request":    "Invalid request data",
				"conflict":       "Data already exists",
				"not_found":      "Resource not found",
				"forbidden":      "You don't have access",
				"unauthorized":   "Unauthorized",
				"internal_error": "Internal server error",
			},
			"id": {
				"bad_request":    "Data request tidak valid",
				"conflict":       "Data sudah ada",
				"not_found":      "Data tidak ditemukan",
				"forbidden":      "Tidak memiliki akses",
				"unauthorized":   "Tidak terotorisasi",
				"internal_error": "Terjadi kesalahan pada server",
			},
		},
	}
}

func (t *Translator) Translate(lang, key string) string {
	lang = normalizeLang(lang)

	if msg, ok := t.messages[lang][key]; ok {
		return msg
	}

	// fallback ke default
	return t.messages[t.defaultLang][key]
}

func normalizeLang(lang string) string {
	if len(lang) == 0 {
		return ""
	}
	lang = strings.ToLower(lang)
	if strings.HasPrefix(lang, "id") {
		return "id"
	}
	if strings.HasPrefix(lang, "en") {
		return "en"
	}
	return "en"
}
