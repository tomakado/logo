package log_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/tomakado/logo/log"
)

func TestLogger_Write(t *testing.T) {
	t.Run("too low logging level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &log.Logger{
			Output:    &buf,
			Level:     log.LevelImportant,
			Formatter: &log.JSONFormatter{},
		}

		logger.Write(context.Background(), log.LevelVerbose, "hello", nil)
		assert.Equal(t, 0, buf.Len())
	})

	t.Run("message logged", func(t *testing.T) {
		tmpl, err := template.New("test_message_logged").Parse("{{.Level}} {{.Message}} {{.Extra}}")
		assert.NoError(t, err)

		formatter := log.NewTemplateFormatter(tmpl)

		var buf bytes.Buffer
		logger := &log.Logger{
			Level:     log.LevelVerbose,
			Output:    &buf,
			Formatter: formatter,
		}

		msg := "hello"
		extra := map[string]interface{}{"foo": "bar"}

		event := log.NewEvent(log.LevelVerbose, msg, extra)
		m := formatter.Format(event)

		logger.Write(context.Background(), log.LevelVerbose, msg, extra)
		assert.Equal(t, len(string(m))+1, len(buf.String()))
	})

	t.Run("empty message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &log.Logger{
			Level:     log.LevelVerbose,
			Output:    &buf,
			Formatter: &log.JSONFormatter{},
		}

		logger.Write(context.Background(), log.LevelVerbose, nil, nil)
		assert.Equal(t, 0, len(buf.String()))
	})
}

func TestLogger_Writef(t *testing.T) {
	var buf bytes.Buffer
	logger := &log.Logger{
		Output:    &buf,
		Formatter: &log.JSONFormatter{},
	}

	const (
		msgFmt = "Hello, %s"
		val    = "Jon Snow"
	)

	logger.Writef(context.Background(), log.LevelVerbose, msgFmt, val)
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf(msgFmt, val)))
}

func TestLogger_Verbose(t *testing.T) {
	t.Run("usual case", func(t *testing.T) {
		var loggedEvent *log.Event
		logger := &log.Logger{
			Output:    ioutil.Discard,
			Formatter: &log.JSONFormatter{},
		}
		logger.PostHook(func(_ context.Context, e *log.Event) {
			loggedEvent = e
		})

		const msg = "hello"

		logger.Verbose(context.Background(), msg)
		assert.Equal(t, log.LevelVerbose, loggedEvent.Level)
		assert.Equal(t, msg, loggedEvent.Message)
	})

	t.Run("too low logging level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := &log.Logger{
			Output:    &buf,
			Level:     log.LevelImportant,
			Formatter: &log.JSONFormatter{},
		}

		logger.Verbose(context.Background(), "hello")
		assert.Equal(t, 0, len(buf.String()))
	})
}

func TestLogger_Important(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Level:     log.LevelImportant,
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	const msg = "hello"

	logger.Important(context.Background(), msg)
	assert.Equal(t, log.LevelImportant, loggedEvent.Level)
	assert.Equal(t, msg, loggedEvent.Message)
}

func TestLogger_VerboseX(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	const msg = "hello"
	extra := log.Extra{"foo": "bar"}

	logger.VerboseX(context.Background(), msg, extra)
	assert.Equal(t, log.LevelVerbose, loggedEvent.Level)
	assert.Equal(t, msg, loggedEvent.Message)
	assert.Equal(t, extra, loggedEvent.Extra)
}

func TestLogger_ImportantX(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Level:     log.LevelImportant,
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	const msg = "hello"
	extra := log.Extra{"foo": "bar"}

	logger.ImportantX(context.Background(), msg, extra)
	assert.Equal(t, log.LevelImportant, loggedEvent.Level)
	assert.Equal(t, msg, loggedEvent.Message)
	assert.Equal(t, extra, loggedEvent.Extra)
}

func TestLogger_Verbosef(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	const (
		msgFmt = "Hello, %s"
		val    = "Jon Snow"
	)

	logger.Verbosef(context.Background(), msgFmt, val)
	assert.Equal(t, log.LevelVerbose, loggedEvent.Level)
	assert.Equal(t, fmt.Sprintf(msgFmt, val), loggedEvent.Message)
}

func TestLogger_Importantf(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Level:     log.LevelImportant,
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PostHook(func(_ context.Context, e *log.Event) {
		loggedEvent = e
	})

	const (
		msgFmt = "Hello, %s"
		val    = "Jon Snow"
	)

	logger.Importantf(context.Background(), msgFmt, val)
	assert.Equal(t, log.LevelImportant, loggedEvent.Level)
	assert.Equal(t, fmt.Sprintf(msgFmt, val), loggedEvent.Message)
}

func TestLogger_PreHook(t *testing.T) {
	var loggedEvent *log.Event
	logger := &log.Logger{
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}
	logger.PreHook(func(_ context.Context, e *log.Event) {
		e.Extra["foo"] = "bar"
	})
	logger.PostHook(func(c context.Context, e *log.Event) {
		loggedEvent = e
	})

	logger.Verbose(context.Background(), "hello")
	assert.Equal(t, "bar", loggedEvent.Extra["foo"])
}

func TestLogger_PostHook(t *testing.T) {
	logger := &log.Logger{
		Level:     log.LevelVerbose,
		Output:    ioutil.Discard,
		Formatter: &log.JSONFormatter{},
	}

	var hookCalled bool
	logger.PostHook(func(_ context.Context, _ *log.Event) {
		hookCalled = true
	})

	logger.Write(context.Background(), log.LevelVerbose, "hello", nil)
	assert.True(t, hookCalled)
}
