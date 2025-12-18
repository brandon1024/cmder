package flag

import (
	"strings"
)

// Flag represents a single short or long flag. It is identified by its name (e.g. 'h' for '-h' and 'help' for
// '--help'), and the value of the flag is represented by the [Value] interface.
type Flag struct {
	// Name is the name of the flag, as it appears at the command line.
	Name string

	// Usage is a descriptive help message for the flag.
	Usage string

	// Value represents the value of the flag as parsed.
	Value Value

	// DefValue is the (stringified) default value for the flag.
	DefValue string
}

// UnquoteUsage extracts a back-quoted name from the usage string for a [Flag] and returns it and the un-quoted usage.
//
// Given
//
//	var output string
//	fs := flag.NewFlagSet("echo", ContinueOnError)
//	fs.StringVar(&output, "output", "-", "output `file` location")
//
// UnquoteUsage would return
//
//	("file", "output file location").
//
// If there are no back quotes, the name is an educated guess of the type of the flag's value, or the empty string if
// the flag is boolean.
func UnquoteUsage(flg *Flag) (string, string) {
	var (
		usage = flg.Usage
		begin = strings.Index(usage, "`")
		end   = strings.LastIndex(usage, "`")
	)

	if begin != -1 && end != -1 {
		usage = flg.Usage[0:begin] + flg.Usage[begin+1:end] + flg.Usage[end+1:]
	}

	if isBoolFlag(flg) {
		return "", usage
	}

	if begin != -1 && end != -1 {
		return flg.Usage[begin+1 : end], usage
	}

	switch flg.Value.(type) {
	case *durationT:
		return "duration", usage
	case *float64T:
		return "float", usage
	case *intT, *int64T:
		return "int", usage
	case *stringT, *textT:
		return "string", usage
	case *uintT, *uint64T:
		return "uint", usage
	default:
		return "value", usage
	}
}

// isBoolFlag checks if the given flag is a boolean flag. Boolean flags do not accept arguments.
func isBoolFlag(flg *Flag) bool {
	bf, ok := flg.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}
