# logo

Opinionated and minimalistic logging library for Go.

Library provides only two levels out of box - *Verbose* and *Important.* Why so? It's mostly inspired by 🇷🇺 [this post](https://t.me/nikitonsky_pub/47) by [@tonsky](https://github.com/tonsky) and my personal experience.

TLDR:
- Only two logging levels: verbose and important.
- *Important* level is for errors and business-critical stuff. *Verbose* level is for development purposes.
- Stuff like "successfully connected to host ABC", "binded port 8000", etc. are not needed even at verbose level.
- Libraries only use important level because debug related things are interesting only for library developers.

However library allows to define custom levels. But before doing it, please think carefully.

# Installation

```bash
go get github.com/tomakado/logo
```

# Usage

You can use pre-instantiated logger and wrapper functions around
it or create and customize your own.

## Quick start

For quick start use package functions like this:

```golang
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
```

For fine-tuned logger use `NewLogger` function:

```golang
package main

import (
    "context"
    "os"

    "github.com/tomakado/logo/log"
)

func main() {
    ctx := context.Background()

    logger := log.NewLogger(log.LevelImportant, os.Stderr, log.SimpleTextFormatter)

    logger.Verbose(ctx, "hello!") // will not be sent to output
    logger.Important(ctx, "this is really important")
}
```

## Logging levels

logo's logging level is a pair of numeric value and string representation of level and can be defined with `NewLevel` function:

```golang
var (
    LevelVerbose   Level = NewLevel(10, "VERBOSE")
    LevelImportant Level = NewLevel(20, "IMPORTANT")
)
```

## Message format

`NewLogger` accepts `Formatter` as third argument to create logger. There are two formatter types out of box `JSONFormatter` and `TemplateFormatter` and two pre-instantiated template formatters: `SimpleTextFormatter` and `TableTextFormatter`.

`TemplateFormatter` uses template engine from Go's standard library to format messages:

```golang
tmpl, err := template.New("tmpl_fmt_example").
    Parse("ts={{.Time}} level={{.Level}} msg=\"{{.Message}}\" extra={{.Extra}}")
if err != nil {
    panic(err)
}

formatter := log.NewTemplateFormatter(tmpl)
logger := log.NewLogger(log.LevelVerbose, os.Stdout, formatter)

logger.Verbose(context.Background(), "hello")
```

## Hooks

Hooks are functions called before or after log message has been sent to output. Pre-hooks are useful when you need to extend the context of event. Post-hooks can be used to send events to external services (e.g. Sentry), collect metrics, etc.

```golang
log.PreHook(func(ctx context.Context, e *log.Event) {
    e.Extra["request_id"] = uuid.New()
})

log.PostHook(func(ctx context.Context, e *log.Event) {
    if e.Level >= log.LevelImportant {
        // For example, send event to Sentry or:
        os.Exit(1)
    }
})
```