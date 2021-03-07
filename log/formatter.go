package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// Formatters to quick start
var (
	TableTextFormatter  *TemplateFormatter
	SimpleTextFormatter *TemplateFormatter
)

// Formatter converts given event to string
type Formatter interface {
	Format(event Event) string
}

// JSONFormatter is used to output logs as JSON string
type JSONFormatter struct{}

// Format converts given event to JSON string
func (f JSONFormatter) Format(event Event) string {
	if event.Message == nil {
		return ""
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// TemplateFormatter is used to output logs rendered with template
type TemplateFormatter struct {
	tmpl *template.Template
}

// NewTextFormatter creates a new instance of TemplateFormatter with given template
func NewTextFormatter(tmpl *template.Template) *TemplateFormatter {
	return &TemplateFormatter{
		tmpl: tmpl,
	}
}

// Format renders event to string with template
func (f TemplateFormatter) Format(event Event) string {
	if event.Message == nil {
		return ""
	}

	var builder strings.Builder
	if err := f.tmpl.Execute(&builder, event); err != nil {
		panic(err)
	}

	return builder.String()
}

func init() {
	createTableTextFormatter()
}

func createTableTextFormatter() {
	tmpl, err := template.New("tablefmt").
		Funcs(template.FuncMap{
			"fmtTime": func(t time.Time) string {
				return t.Format(time.RubyDate)
			},
			"fmtLevel": func(level Level) string {
				return fmt.Sprintf("%-9s", level)
			},
		}).
		Parse("| {{.Level|fmtLevel}} | {{.Time|fmtTime}} | {{.Message}}{{if .Extra}}; extra: {{.Extra}} {{end}}")
	if err != nil {
		panic(err)
	}

	TableTextFormatter = NewTextFormatter(tmpl)
}

func createSimpleTextFormatter() {
	tmpl, err := template.New("simplefmt").Parse("{{ .Level}} @ {{.Time}}: {{.Message}}{{if .Extra}}; {{.Extra}}")
	if err != nil {
		panic(err)
	}

	SimpleTextFormatter = NewTextFormatter(tmpl)
}
