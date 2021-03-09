package hooks

import (
	"context"

	"github.com/google/uuid"
	"github.com/tomakado/logo/log"
)

// EventID adds unique identifier to each log event.
func EventID(_ context.Context, e *log.Event) {
	e.Extra["event_id"] = uuid.New()
}
