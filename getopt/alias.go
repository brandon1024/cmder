package getopt

import (
	"flag"
	"fmt"
)

// Alias is a simply utility for registering flag aliases, registering a new flag in fs with name alias with the
// [flag.Value] of a flag named name.
//
// If flag name doesn't exist in fs, panic.
func Alias(fs *flag.FlagSet, name, alias string) {
	flg := fs.Lookup(name)
	if flg == nil {
		panic(fmt.Sprintf("cmder: cannot register alias '%s: target '%s' does not exist in flag set'", alias, name))
	}

	fs.Var(flg.Value, alias, flg.Usage)
}
