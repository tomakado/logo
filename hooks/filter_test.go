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

func TestFilteredHook(t *testing.T) {
	var hookCalled bool

	hook := func(_ context.Context, _ *log.Event) {
		hookCalled = true
	}

	// this filter does not call hook on events without "request_id" key in extra
	filter := func(e *log.Event) bool {
		if _, ok := e.Extra["request_id"]; !ok {
			return false
		}

		return true
	}

	filteredHook := hooks.FilteredHook(hook, filter)

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
	logger.PostHook(filteredHook)

	ctx := context.Background()

	logger.Verbose(ctx, "hook will not be called on this message")
	assert.False(t, hookCalled)

	logger.VerboseX(ctx, "hook will be called on this message", log.Extra{"request_id": uuid.New()})
	assert.True(t, hookCalled)
}

func TestLevelBoundsFilter(t *testing.T) {
	var hookCalled bool

	hook := func(_ context.Context, _ *log.Event) {
		hookCalled = true
	}

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
	logger.PostHook(
		hooks.FilteredHook(
			hook,
			hooks.LevelBoundsFilter(log.LevelVerbose, log.LevelImportant),
		),
	)

	outLevel := log.NewLevel(100, "MOST_CRITICAL_EVER")
	ctx := context.Background()

	logger.Write(ctx, outLevel, "very critical stuff", nil)
	assert.False(t, hookCalled)

	logger.Important(ctx, "hello world")
	assert.True(t, hookCalled)
}

func TestLevelFilter(t *testing.T) {
	var hookCalled bool

	hook := func(_ context.Context, _ *log.Event) {
		hookCalled = true
	}

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
	logger.PostHook(
		hooks.FilteredHook(
			hook,
			hooks.LevelFilter(log.LevelImportant),
		),
	)

	ctx := context.Background()

	logger.Verbose(ctx, "not so important")
	assert.False(t, hookCalled)

	logger.Important(ctx, "really important")
	assert.True(t, hookCalled)
}
