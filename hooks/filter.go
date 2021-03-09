package hooks

import (
	"context"

	"github.com/tomakado/logo/log"
)

func Filter(h log.Hook, minLevel, maxLevel log.Level) log.Hook {
	return func(ctx context.Context, e *log.Event) {
		if !levelHitsBounds(e.Level, minLevel, maxLevel) {
			return
		}

		h(ctx, e)
	}
}

func levelHitsBounds(level, minLevel, maxLevel log.Level) bool {
	return level.Gte(minLevel) && maxLevel.Gte(level)
}
