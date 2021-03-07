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

	t.Run("hooks called", func(t *testing.T) {
		logger := &log.Logger{
			Level:     log.LevelVerbose,
			Output:    ioutil.Discard,
			Formatter: &log.JSONFormatter{},
		}

		var hookCalled bool
		logger.Hook(func(_ context.Context, _ log.Event) {
			hookCalled = true
		})

		logger.Write(context.Background(), log.LevelVerbose, "hello", nil)

		assert.True(t, hookCalled)
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
		var loggedEvent log.Event
		logger := &log.Logger{
			Output:    ioutil.Discard,
			Formatter: &log.JSONFormatter{},
		}
		logger.Hook(func(_ context.Context, e log.Event) {
			loggedEvent = e
		})

		const msg = "hello"

		logger.Verbose(context.Background(), msg)
		assert.Equal(t, log.LevelVerbose, loggedEvent.Level)
		assert.Equal(t, loggedEvent.Message, msg)
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