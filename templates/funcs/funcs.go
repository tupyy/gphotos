package funcs

import "time"

func Day(t time.Time) int {
	return int(t.Day())
}
