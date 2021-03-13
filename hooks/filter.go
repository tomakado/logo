package hooks

import (
	"context"

	"github.com/tomakado/logo/log"
)

// FilteredHook returns log.Hook with filter before original function.
func FilteredHook(h log.Hook, filter Filter) log.Hook {
	return func(ctx context.Context, e *log.Event) {
		if !filter(e) {
			return
		}

		h(ctx, e)
	}
}

// Filter is a function that returns true in case hook should be called
// and false otherwise.
type Filter func(e *log.Event) bool

// LevelBoundsFilter returns filter based on given logging level bounds.
func LevelBoundsFilter(min, max log.Level) Filter {
	return func(e *log.Event) bool {
		return e.Level.Gte(min) && max.Gte(e.Level)
	}
}

// LevelFilter returns filter that checks if event level is greater or
// equal to filter's logging level.
func LevelFilter(level log.Level) Filter {
	return func(e *log.Event) bool {
		return e.Level.Gte(level)
	}
}
