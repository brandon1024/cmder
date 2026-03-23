package getopt

import (
	"flag"
	"fmt"
)

// Alias is a simple utility for registering flag aliases. A new flag is registered in fs with name alias and the
// [flag.Value] of a flag named name.
//
// If flag name doesn't exist in fs, panic.
func Alias(fs *flag.FlagSet, name, alias string) {
	flg := fs.Lookup(name)
	if flg == nil {
		panic(fmt.Sprintf("getopt: cannot register alias '%s': target '%s' does not exist in flag set", alias, name))
	}

	fs.Var(flg.Value, alias, flg.Usage)
}
