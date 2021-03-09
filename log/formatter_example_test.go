package log_test

import (
	"context"
	"os"
	"text/template"

	"github.com/tomakado/logo/log"
)

func ExampleTemplateFormatter() {
	tmpl, err := template.New("tmpl_fmt_example").
		Parse("level={{.Level}} msg=\"{{.Message}}\" extra={{.Extra}}")
	if err != nil {
		panic(err)
	}

	formatter := log.NewTemplateFormatter(tmpl)
	logger := log.NewLogger(log.LevelVerbose, os.Stdout, formatter)

	logger.Verbose(context.Background(), "hello")
	// Output: level=VERBOSE msg="hello" extra=map[]
}
