package getopt

import "flag"

// HiddenFlag is a [flag.Value] that should be hidden from [PosixFlagSet.PrintDefaults] output.
type HiddenFlag interface {
	flag.Value
	IsHiddenFlag() bool
}

// Hidden is a [flag.Value] that is hidden from [PosixFlagSet.PrintDefaults] output.
type Hidden struct {
	flag.Value
}

// Hide hides the given flag.
func Hide(flg *flag.Flag) {
	flg.Value = &Hidden{flg.Value}
}

// IsHiddenFlag implements [HiddenFlag] and returns true.
func (h *Hidden) IsHiddenFlag() bool {
	return true
}

// String returns the parent [flag.Value].
func (h *Hidden) String() string {
	// if [Hidden] is used with the standard [flag.FlagSet], its [PrintDefaults] will call this method on a zero value,
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
