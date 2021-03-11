package log_test

import (
	"encoding/json"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomakado/logo/log"
)

func TestJSONFormatter_Format(t *testing.T) {
	formatter := &log.JSONFormatter{}

	t.Run("empty message", func(t *testing.T) {
		event := log.NewEvent(log.LevelVerbose, nil, map[string]interface{}{"foo": "bar"})

		assert.Equal(t, "", formatter.Format(event))
	})

	t.Run("usual case", func(t *testing.T) {
		eventTime := time.Now()
		event := log.Event{
			Time:    eventTime,
			Level:   log.LevelVerbose,
			Message: "hello",
			Extra:   map[string]interface{}{"foo": "bar"},
		}

		jsonEvent := struct {
			Level string `json:"level"`
			*log.Event
		}{
			Level: event.Level.String(),
			Event: &event,
		}

		m, err := json.Marshal(jsonEvent)
		assert.NoError(t, err)

		assert.Equal(t, string(m), formatter.Format(event))
	})
}

func TestTemplateFormatter(t *testing.T) {
	t.Run("empty message", func(t *testing.T) {
		tmpl, err := template.New("test_empty_message").Parse("")
		assert.NoError(t, err)

		formatter := log.NewTemplateFormatter(tmpl)

		event := log.NewEvent(log.LevelVerbose, nil, map[string]interface{}{"foo": "bar"})

		assert.Equal(t, "", formatter.Format(event))
	})

	t.Run("usual case", func(t *testing.T) {
		tmpl, err := template.New("test_usual_case").Parse("{{.Level}} {{.Time}} {{.Message}} {{.Extra}}")
		assert.NoError(t, err)

		formatter := log.NewTemplateFormatter(tmpl)

		eventTime := time.Now()
		event := log.Event{
			Time:    eventTime,
			Level:   log.LevelVerbose,
			Message: "hello",
			Extra:   map[string]interface{}{"foo": "bar"},
		}

		var rendered strings.Builder
		assert.NoError(t, tmpl.Execute(&rendered, event))

		assert.Equal(t, rendered.String(), formatter.Format(event))
	})
}
