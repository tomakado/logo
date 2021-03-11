package hooks_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tomakado/logo/hooks"
	"github.com/tomakado/logo/log"
)

func TestEventID(t *testing.T) {
	var loggedEvent *log.Event

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
	logger.PreHook(hooks.EventID)
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	logger.Verbose(context.Background(), "hello")
	_, isUUID := (loggedEvent.Extra["event_id"]).(uuid.UUID)
	assert.True(t, isUUID)
}
