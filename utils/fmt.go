package utils

import (
	"fmt"
	"time"
)

func FormatTimeRuby(t time.Time) string {
	return t.Format(time.RubyDate)
}

func FormatLevelFixedWidth(l string) string {
	return fmt.Sprintf("%-9s", l)
}
