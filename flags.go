package cmder

import (
	"flag"
	"reflect"
	"slices"

	"github.com/brandon1024/cmder/getopt"
)

// FlagInitializer is an interface implemented by commands that need to register flags.
//
// InitializeFlags will be invoked during [Execute], prior to Initialize/Run/Destroy routines. You can use this to
// register flags for your command.
//
// If the command does not define help flags '-h' and '--help', they will be registered automatically and will instruct
// [Execute] to render command usage.
type FlagInitializer interface {
	InitializeFlags(*flag.FlagSet)
}

// flagParser is an interface implemented by types that parse args.
type flagParser interface {
	Parse([]string) error
	Args() []string
}

// areSame check if f1 and f2 have the same underlying [flag.Value].
func areSame(f1, f2 flag.Value) bool {
	var (
		ref1 = reflect.ValueOf(f1)
		ref2 = reflect.ValueOf(f2)
	)

	if ref1.Comparable() && ref2.Comparable() && f1 == f2 {
		return true
	}

	if ref1.Kind() != ref2.Kind() {
		return false
	}

	if !slices.Contains([]reflect.Kind{reflect.Map, reflect.Pointer, reflect.Func, reflect.Slice}, ref1.Kind()) {
		return false
	}

	return ref1.Pointer() == ref2.Pointer()
}

// isHiddenFlag checks if the given flag is hidden.
func isHiddenFlag(flg *flag.Flag) bool {
	hf, ok := flg.Value.(getopt.HiddenFlag)
	return ok && hf.IsHiddenFlag()
}

// boolFlag is a [flag.Value] that implements an additional method IsBoolFlag which indicates whether the flag accepts
// arguments or not.
type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// isBoolFlag checks if the given flag is a boolean flag.
func isBoolFlag(flg *flag.Flag) bool {
	hf, ok := flg.Value.(boolFlag)
	return ok && hf.IsBoolFlag()
}
