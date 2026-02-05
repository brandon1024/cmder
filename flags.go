package cmder

import (
	"flag"
)

// FlagInitializer is an interface implemented by commands that need to register flags.
//
// InitializeFlags will be invoked during [Execute], prior to any lifecycle routines. You can use this to register
// flags for your command.
//
// Help flags '-h' and '--help' are registered automatically and will instruct [Execute] to render usage information
// to the [UsageOutputWriter].
type FlagInitializer interface {
	InitializeFlags(*flag.FlagSet)
}
