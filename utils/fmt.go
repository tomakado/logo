package utils

import (
	"fmt"
	"time"
)

// FormatTimeRuby returns time.Time value formatted in Ruby style.
func FormatTimeRuby(t time.Time) string {
	return t.Format(time.RubyDate)
}

// FormatLevelFixedWidth returns given level string representation with constant width.
func FormatLevelFixedWidth(l string) string {
	return fmt.Sprintf("%-9s", l)
}
