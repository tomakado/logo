package main

import (
	"context"
	"os"
	"text/template"

	"github.com/tomakado/logo/log"
)

type someInfo struct {
	A string
	B int
}

func main() {
	logger := log.DefaultLogger
	ctx := context.Background()

	logger.Write(ctx, log.LevelImportant, "hello!", map[string]string{"foo": "bar"})
	logger.VerboseX(ctx, "hello again", someInfo{A: "foo", B: 17})
	logger.Important(ctx, "this is important info!")

	tmpl, err := template.New("fmt").Parse("{{.Time|toStamp}} | {{.Level}} | {{.Message}} | {{.Extra}}")
	if err != nil {
		logger.Important(ctx, err)
	}
	fmtr := log.NewTextFormatter(tmpl)

	tmplLogger := &log.Logger{
		Level:     log.LevelImportant,
		Output:    os.Stderr,
		Formatter: fmtr,
	}

	tmplLogger.Important(ctx, "wow!")
	tmplLogger.Verbose(ctx, "you will not see this")
}
