package i18n

import (
	"path"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tupyy/gophoto/internal/conf"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// load languages
	Bundle.MustLoadMessageFile(path.Join(conf.GetStaticsFolder(), "i18n/active.en.toml"))
	Bundle.MustLoadMessageFile(path.Join(conf.GetStaticsFolder(), "i18n/active.ro.toml"))
}

func GetTranslation(localizer *i18n.Localizer, id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}
