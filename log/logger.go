package log

import (
	"context"
	"io"
	"os"
	"strings"
)

var DefaultLogger *Logger

type Logger struct {
	Level     Level
	Output    io.Writer
	Formatter Formatter
	Hooks     []Hook
}

type Hook func(context.Context, Event)

func (l *Logger) Verbose(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelVerbose, msg, nil)
}

func (l *Logger) Important(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelImportant, msg, nil)
}

func (l *Logger) VerboseX(ctx context.Context, msg interface{}, extra interface{}) {
	l.Write(ctx, LevelVerbose, msg, extra)
}

func (l *Logger) ImportantX(ctx context.Context, msg interface{}, extra interface{}) {
	l.Write(ctx, LevelImportant, msg, extra)
}

func (l *Logger) Write(ctx context.Context, level Level, msg interface{}, extra interface{}) {
	event := NewEvent(level, msg, extra)
	formattedEvent := l.Formatter.Format(event) + "\n"

	l.Output.Write([]byte(formattedEvent))

	for _, h := range l.Hooks {
		h(ctx, event)
	}
}

func (l *Logger) Hook(h Hook) {
	l.Hooks = append(l.Hooks, h)
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
