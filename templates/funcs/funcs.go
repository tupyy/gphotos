package funcs

import (
	"fmt"
	"strings"
	"time"

	goI18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/i18n"
)

func Day(t time.Time) int {
	return int(t.Day())
}

func Month(t time.Time) string {
	return t.Month().String()[:3]
}

func Year(t time.Time) int {
	return t.Year()
}

func PermissionName(p entity.Permission) string {
	str := strings.Split(p.String(), ".")
	return string(str[1][0])
}

func Date(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// set a nicer format of exif format
func DatePhoto(s string) string {
	format := "2006:01:02 15:04:05"
	outputFormat := "02/01/2006 15:04:05"
	if t, err := time.Parse(format, s); err == nil {
		return t.Format(outputFormat)
	}

	return s
}

func ExtractMetadata(name string, metadata map[string]string) string {
	key := fmt.Sprintf("X-Amz-Meta-%s", strings.Title(name))
	if value, found := metadata[key]; found {
		return value
	}
	return "N/A"
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToTitle(s string) string {
	return strings.Title(s)
}

func CapFirst(s string) string {
	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func Translate(lang, id string) string {
	localizer := goI18n.NewLocalizer(i18n.Bundle, lang)
	return i18n.GetTranslation(localizer, id)
}
