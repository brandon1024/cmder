package cmder

import (
	"flag"
)

// Implemented by commands that need to register flags.
type FlagInitializer interface {
	// InitializeFlags initializes flags. Invoked by [Execute] before any lifecycle routines.
	InitializeFlags(*flag.FlagSet)
}
