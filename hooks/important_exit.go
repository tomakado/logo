package hooks

import (
	"context"
	"os"

	"github.com/tomakado/logo/log"
)

var (
	// ExitOnImportant exits from program with code 1 if event level is greater than or equal to important.
	ExitOnImportant = FilteredHook(
		func(_ context.Context, _ *log.Event) { os.Exit(1) },
		LevelFilter(log.LevelImportant),
	)
)
