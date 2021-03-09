package hooks

import (
	"context"
	"os"

	"github.com/tomakado/logo/log"
)

// ExitOnImportant exits from program with code 1 if event level is greater than or equal to important.
func ExitOnImportant(_ context.Context, e *log.Event) log.Hook {
	return Filter(
		func(_ context.Context, _ *log.Event) { os.Exit(1) },
		log.LevelImportant,
		log.LevelImportant,
	)
}
