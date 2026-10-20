package getopt

import "flag"

// boolFlag is a [flag.Value] that also implements a method IsBoolFlag, used to determine if the flag accepts an
// argument or not.
type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// isBoolFlag checks if the given flag has a [flag.Value] which is a boolean flag.
func isBoolFlag(flg *flag.Flag) bool {
	bf, ok := flg.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}
