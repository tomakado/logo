package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// DefaultLogger for quick start
var DefaultLogger *Logger

// Logger ...
type Logger struct {
	Level     Level
	Output    io.Writer
	Formatter Formatter
	Hooks     []Hook
}

// Hook is a function being called when event was sent to logger output
type Hook func(context.Context, Event)

// Verbose writes a message with verbose level
func (l *Logger) Verbose(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelVerbose, msg, nil)
}

// Important writes message with important level
func (l *Logger) Important(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelImportant, msg, nil)
}

// VerboseX writes a message with verbose level and given extra
func (l *Logger) VerboseX(ctx context.Context, msg interface{}, extra interface{}) {
	l.Write(ctx, LevelVerbose, msg, extra)
}

// ImportantX writes a message with important level and given extra
func (l *Logger) ImportantX(ctx context.Context, msg interface{}, extra interface{}) {
	l.Write(ctx, LevelImportant, msg, extra)
}

// Verbosef writes a formatted message with verbose level
func (l *Logger) Verbosef(ctx context.Context, msg string, values ...interface{}) {
	l.Writef(ctx, LevelVerbose, msg, values...)
}

// Importantf writes a formatted message with important level
func (l *Logger) Importantf(ctx context.Context, msg string, values ...interface{}) {
	l.Writef(ctx, LevelImportant, msg, values...)
}

// Writef writes a formatted message with given level
func (l *Logger) Writef(ctx context.Context, level Level, msg string, values ...interface{}) {
	l.Write(ctx, level, fmt.Sprintf(msg, values...), nil)
}

// Write writes a message with given level and exra
func (l *Logger) Write(ctx context.Context, level Level, msg interface{}, extra interface{}) {
	if l.Level == LevelImportant && level == LevelVerbose {
		return
	}

	if msg == nil {
		return
	}

	event := NewEvent(level, msg, extra)
	formattedEvent := l.Formatter.Format(event) + "\n"

	l.Output.Write([]byte(formattedEvent))

	for _, h := range l.Hooks {
		h(ctx, event)
	}
}

// Hook register given hooks in logger
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
