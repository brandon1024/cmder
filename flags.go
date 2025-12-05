package cmder

import (
	"flag"
)

// FlagInitializer is an interface implemented by commands that need to register flags.
type FlagInitializer interface {
	// InitializeFlags initializes flags. Invoked by [Execute] before any lifecycle routines.
	//
	// Help flags '-h' and '--help' are registered automatically and will instruct [Execute] to render usage information
	// to the [UsageOutputWriter].
	//
	// See [Execute] for more information.
	InitializeFlags(*flag.FlagSet)
}
