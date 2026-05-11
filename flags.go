package cmder

import (
	"flag"
)

// FlagInitializer is an interface implemented by a [Command] that need to register flags.
//
// InitializeFlags will be invoked during [Execute], prior to Initialize()/Run()/Destroy() routines. You can use this to
// register flags for your command.
//
// If the command does not define help flags '-h' or '--help', they will be registered automatically and will instruct
// [Execute] to render command usage.
type FlagInitializer interface {
	InitializeFlags(*flag.FlagSet)
}

// flagParser is an interface implemented by types that parse args.
type flagParser interface {
	Parse([]string) error
	Args() []string
}
