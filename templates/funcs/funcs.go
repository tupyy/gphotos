package funcs

import (
	"strings"
	"time"

	"github.com/tupyy/gophoto/internal/entity"
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

func TagName(tag []string) string {
	return tag[0]
}

func TagColor(tag []string) string {
	if len(tag) == 2 && tag[1] != "" {
		return tag[1]
	}

	return "black"
}
