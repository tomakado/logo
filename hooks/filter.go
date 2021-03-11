package hooks

import (
	"context"

	"github.com/tomakado/logo/log"
)

func FilteredHook(h log.Hook, filter Filter) log.Hook {
	return func(ctx context.Context, e *log.Event) {
		if !filter(e) {
			return
		}

		h(ctx, e)
	}
}

type Filter func(e *log.Event) bool

func LevelBoundsFilter(min, max log.Level) Filter {
	return func(e *log.Event) bool {
		return e.Level.Gte(min) && max.Gte(e.Level)
	}
}

func LevelFilter(level log.Level) Filter {
	return func(e *log.Event) bool {
		return e.Level.Gte(level)
	}
}
