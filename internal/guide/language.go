package guide

import (
	"fmt"
	"strings"
)

type Locale string

const (
	LocaleEnglish  Locale = "en"
	LocaleChinese  Locale = "zh"
	LocaleJapanese Locale = "ja"
)

type Translator struct {
	Locale  Locale
	phrases map[string]string
}

func NewTranslator(locale Locale) Translator {
	phrases := map[string]string{}
	switch locale {
	case LocaleChinese:
		phrases = map[string]string{"meet": "集合", "stops": "站点", "confirmed": "已确认", "guests": "游客"}
	case LocaleJapanese:
		phrases = map[string]string{"meet": "集合場所", "stops": "立ち寄り", "confirmed": "確認済み", "guests": "ゲスト"}
	default:
		phrases = map[string]string{"meet": "Meet", "stops": "Stops", "confirmed": "Confirmed", "guests": "Guests"}
	}
	return Translator{Locale: locale, phrases: phrases}
}

func (t Translator) Phrase(key string) string {
	if value, ok := t.phrases[key]; ok {
		return value
	}
	return key
}

func (t Translator) Label(key, value string) string {
	return fmt.Sprintf("%s: %s", t.Phrase(key), value)
}

func (t Translator) ConfirmationLabels(route, meeting string, party int) []string {
	return []string{t.Label("confirmed", route), t.Label("meet", meeting), t.Label("guests", fmt.Sprintf("%d", party))}
}

func SupportedLocales() []Locale { return []Locale{LocaleEnglish, LocaleChinese, LocaleJapanese} }

func ParseLocale(value string) Locale {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, locale := range SupportedLocales() {
		if string(locale) == value {
			return locale
		}
	}
	return LocaleEnglish
}

func IsRightToLeft(locale Locale) bool { return false }

func JoinLabels(labels []string) string { return strings.Join(labels, " | ") }
