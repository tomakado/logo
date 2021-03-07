package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

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
		var buf bytes.Buffer
		logger := &log.Logger{
			Level:     log.LevelVerbose,
			Output:    &buf,
			Formatter: &log.JSONFormatter{},
		}

		msg := "hello"
		extra := map[string]interface{}{"foo": "bar"}

		event := log.NewEvent(log.LevelVerbose, msg, extra)
		m, err := json.Marshal(event)
		assert.NoError(t, err)

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

	msgFmt := "Hello, %s"
	val := "Jon Snow"

	logger.Writef(context.Background(), log.LevelVerbose, msgFmt, val)

	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf(msgFmt, val)))
}
