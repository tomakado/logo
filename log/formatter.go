package log

import (
	"encoding/json"
	"strings"
	"text/template"
)

type Formatter interface {
	Format(event Event) string
}

type JSONFormatter struct {
}

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

type TemplateFormatter struct {
	tmpl *template.Template
}

func NewTextFormatter(tmpl *template.Template) TemplateFormatter {
	return TemplateFormatter{
		tmpl: tmpl,
	}
}

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
