/*
Package log implements flexible but simple logging toolchain.
You can use pre-instantiated logger and wrapper functions around
it or create and customize your own.

For quick start use package functions like this:
	package main

	import (
		"context"

		"github.com/tomakado/logo/log"
	)

	func main() {
		ctx := context.Background()

		log.Verbose(ctx, "hello!")
		log.Important(ctx, "hello, it's important")
		log.VerboseX(ctx, "hello with extra!", log.Extra{"foo": "bar"})
		log.Verbosef(ctx, "hello, %s", "Jon Snow")

		log.Write(ctx, log.LevelImportant, "hello, it's me", Extra{"a": 42})
		log.Writef(ctx, log.LevelVerbose, "My name is %s, I'm %d y.o.", "Ildar", 23)
	}
*/
package log

import (
	"context"
	"os"
	"strings"
)

// DefaultLogger for quick start
var DefaultLogger *Logger

// Verbose writes a message with verbose level
func Verbose(ctx context.Context, msg interface{}) {
	DefaultLogger.Verbose(ctx, msg)
}

// Important writes a message with important level
func Important(ctx context.Context, msg interface{}) {
	DefaultLogger.Important(ctx, msg)
}

// VerboseX writes a message with verbose level and given extra
func VerboseX(ctx context.Context, msg interface{}, extra Extra) {
	DefaultLogger.VerboseX(ctx, msg, extra)
}

// ImportantX writes a message with important level and given extra
func ImportantX(ctx context.Context, msg interface{}, extra Extra) {
	DefaultLogger.ImportantX(ctx, msg, extra)
}

// Verbosef writes a formatted message with verbose level
func Verbosef(ctx context.Context, msg string, values ...interface{}) {
	DefaultLogger.Verbosef(ctx, msg, values...)
}

// Importantf writes a formatted message with important level
func Importantf(ctx context.Context, msg string, values ...interface{}) {
	DefaultLogger.Importantf(ctx, msg, values...)
}

// Writef writes a formatted message with given level
func Writef(ctx context.Context, level Level, msg string, values ...interface{}) {
	DefaultLogger.Writef(ctx, level, msg, values...)
}

// Write writes a message with given level and extra
func Write(ctx context.Context, level Level, msg interface{}, extra Extra) {
	DefaultLogger.Write(ctx, level, msg, extra)
}

func init() {
	lvlEnv := strings.ToLower(os.Getenv("LOG_LEVEL"))
	var lvl Level

	switch lvlEnv {
	case "critical", "important", "fatal", "error", "warning", "err", "warn":
		lvl = LevelImportant
	default:
		lvl = LevelVerbose
	}

	DefaultLogger = &Logger{
		Level:     lvl,
		Output:    os.Stderr,
		Formatter: &JSONFormatter{},
	}
}
