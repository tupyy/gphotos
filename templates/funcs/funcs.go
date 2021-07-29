package funcs

import "time"

func Day(t time.Time) int {
	return int(t.Day())
}

func Month(t time.Time) string {
	return t.Month().String()[:3]
}

func Year(t time.Time) int {
	return t.Year()
}
