package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// load languages
	Bundle.MustLoadMessageFile("i18n/active.en.toml")
	Bundle.MustLoadMessageFile("i18n/active.ro.toml")
}

func GetTranslation(localizer *i18n.Localizer, id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}
