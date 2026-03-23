package getopt

import (
	"flag"
	"fmt"
)

// HiddenFlag is a [flag.Value] that should be hidden from [PosixFlagSet.PrintDefaults] output.
type HiddenFlag interface {
	flag.Value
	IsHiddenFlag() bool
}

// HiddenVar is a [flag.Value] that is hidden from [PosixFlagSet.PrintDefaults] output.
type HiddenVar struct {
	flag.Value
}

// Hide is a simple utility for marking a particular flag as hidden from [PosixFlagSet.PrintDefaults] output. The flag
// [flag.Value] for a named flag in fs will be wrapped with [HiddenVar], signaling that the flag is hidden. This is
// functionally equivalent to:
//
//	flg := fs.Lookup(name)
//	flg.Value = &getopt.HiddenVar{flg.Value}
//
// If flag name doesn't exist in fs, panic.
func Hide(fs *flag.FlagSet, name string) {
	flg := fs.Lookup(name)
	if flg == nil {
		panic(fmt.Sprintf("cmder: cannot hide flag '%s': flag '%s' does not exist in flag set", name, name))
	}

	flg.Value = &HiddenVar{flg.Value}
}

// IsHiddenFlag implements [HiddenFlag] and returns true.
func (h *HiddenVar) IsHiddenFlag() bool {
	return true
}

// String returns the parent [flag.Value].
func (h *HiddenVar) String() string {
	// if [HiddenVar] is used with the standard [flag.FlagSet], its [PrintDefaults] will call this method on a zero value,
	// so check the receiver to avoid panics
	if h == nil || h.Value == nil {
		return ""
	}

	return h.Value.String()
}

// isHiddenFlag checks if the given flag has a [flag.Value] which indicates that flg is hidden.
func isHiddenFlag(flg *flag.Flag) bool {
	hf, ok := flg.Value.(HiddenFlag)
	return ok && hf.IsHiddenFlag()
}
