package log

import (
	"context"
	"fmt"
	"io"
)

// Logger ...
type Logger struct {
	Level     Level
	Output    io.Writer
	Formatter Formatter
	preHooks  []Hook
	postHooks []Hook
}

// Hook is a function being called before event was sent to logger ou
type Hook func(context.Context, *Event)

// Verbose writes a message with verbose level
func (l *Logger) Verbose(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelVerbose, msg, nil)
}

// Important writes a message with important level
func (l *Logger) Important(ctx context.Context, msg interface{}) {
	l.Write(ctx, LevelImportant, msg, nil)
}

// VerboseX writes a message with verbose level and given extra
func (l *Logger) VerboseX(ctx context.Context, msg interface{}, extra Extra) {
	l.Write(ctx, LevelVerbose, msg, extra)
}

// ImportantX writes a message with important level and given extra
func (l *Logger) ImportantX(ctx context.Context, msg interface{}, extra Extra) {
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

// Write writes a message with given level and extra
func (l *Logger) Write(ctx context.Context, level Level, msg interface{}, extra Extra) {
	if msg == nil {
		return
	}

	event := NewEvent(level, msg, extra)
	for _, h := range l.preHooks {
		h(ctx, &event)
	}

	if l.Level.IsHigherThan(level) {
		return
	}

	formattedEvent := l.Formatter.Format(event) + "\n"

	l.Output.Write([]byte(formattedEvent))

	for _, h := range l.postHooks {
		h(ctx, &event)
	}
}

// PreHook registers given hook in logger to be executed before log event was written to output
func (l *Logger) PreHook(h Hook) {
	l.preHooks = append(l.preHooks, h)
}

// PostHook registers given hook in  logger to be executed after log event was written to output
func (l *Logger) PostHook(h Hook) {
	l.postHooks = append(l.postHooks, h)
}
