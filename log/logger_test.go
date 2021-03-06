package log_test

import (
	"bytes"
	"context"
	"errors"
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

		logger := log.NewLogger(log.LevelImportant, &buf, &log.JSONFormatter{})

		logger.Write(context.Background(), log.LevelVerbose, "hello", nil)
		assert.Equal(t, 0, buf.Len())
	})

	t.Run("message logged", func(t *testing.T) {
		tmpl, err := template.New("test_message_logged").Parse("{{.Level}} {{.Message}} {{.Extra}}")
		assert.NoError(t, err)

		formatter := log.NewTemplateFormatter(tmpl)

		var buf bytes.Buffer
		logger := log.NewLogger(log.LevelVerbose, &buf, formatter)

		msg := "hello"
		extra := map[string]interface{}{"foo": "bar"}

		event := log.NewEvent(log.LevelVerbose, msg, extra)
		m, err := formatter.Format(event)
		assert.NoError(t, err)

		logger.Write(context.Background(), log.LevelVerbose, msg, extra)
		assert.Equal(t, len(string(m))+1, len(buf.String()))
	})

	t.Run("empty message", func(t *testing.T) {
		var buf bytes.Buffer

		logger := log.NewLogger(log.LevelVerbose, &buf, &log.JSONFormatter{})

		logger.Write(context.Background(), log.LevelVerbose, nil, nil)
		assert.Equal(t, 0, len(buf.String()))
	})

	t.Run("panic on formatting error", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()

		logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, errorFormatter{})
		logger.Verbose(context.Background(), "hello")
	})

	t.Run("panic on event writing error", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()

		logger := log.NewLogger(log.LevelVerbose, errorWriter{}, &log.JSONFormatter{})
		logger.Verbose(context.Background(), "hello")
	})
}

type errorFormatter struct{}

func (f errorFormatter) Format(e log.Event) (string, error) {
	return "", errors.New("error!")
}

type errorWriter struct{}

func (w errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error!")
}

func TestLogger_Writef(t *testing.T) {
	var buf bytes.Buffer

	logger := log.NewLogger(log.LevelVerbose, &buf, &log.JSONFormatter{})

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

		logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
		logger.PostHook(func(_ context.Context, e *log.Event) {
			loggedEvent = e
		})

		msg := struct {
			Greeting string `json:"greeting"`
		}{
			Greeting: "hello",
		}

		logger.Verbose(context.Background(), msg)
		assert.Equal(t, log.LevelVerbose, loggedEvent.Level)
		assert.Equal(t, msg, loggedEvent.Message)
	})

	t.Run("too low logging level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.NewLogger(log.LevelImportant, &buf, &log.JSONFormatter{})

		logger.Verbose(context.Background(), "hello")
		assert.Equal(t, 0, len(buf.String()))
	})
}

func TestLogger_Important(t *testing.T) {
	var loggedEvent *log.Event

	logger := log.NewLogger(log.LevelImportant, ioutil.Discard, &log.JSONFormatter{})
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

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
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

	logger := log.NewLogger(log.LevelImportant, ioutil.Discard, &log.JSONFormatter{})
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

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
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

	logger := log.NewLogger(log.LevelImportant, ioutil.Discard, &log.JSONFormatter{})
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

	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})
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
	logger := log.NewLogger(log.LevelVerbose, ioutil.Discard, &log.JSONFormatter{})

	var hookCalled bool
	logger.PostHook(func(_ context.Context, _ *log.Event) {
		hookCalled = true
	})

	logger.Write(context.Background(), log.LevelVerbose, "hello", nil)
	assert.True(t, hookCalled)
}
