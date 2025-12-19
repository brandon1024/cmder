package flag

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"time"
	"unicode"
)

var (
	ErrHelp = errors.New("flag: help requested")
)

// CommandLine is the default set of command-line flags parsed from [os.Args].
var (
	CommandLine *FlagSet
)

func init() {
	if len(os.Args) == 0 {
		CommandLine = NewFlagSet("", ExitOnError)
	} else {
		CommandLine = NewFlagSet(os.Args[0], ExitOnError)
	}

	CommandLine.Usage = Usage
}

// Usage is a simple usage function which prints usage information for the global [CommandLine].
var Usage = func() {
	_, err := fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", CommandLine.Name())
	if err != nil {
		panic(err)
	}

	CommandLine.PrintDefaults()
}

// FlagSet is a set of flags. The zero value of a FlagSet has no name and has [ContinueOnError] error handling policy.
//
// Unlike FlagSet in the standard library package, this FlagSet parses flags with POSIX/GNU semantics.
//
// [Flag] names must be unique within a FlagSet. An attempt to define a flag whose name is already in use will cause a
// panic.
type FlagSet struct {
	// Usage is a function called when an error occurs while parsing flags. It is invoked directly after an error is
	// encountered, but immediately before [FlagSet.Parse] returns the error or exits/panics (see [ErrorHandling]).
	//
	// If nil, defaults to [PrintDefaults].
	Usage func()

	name          string
	errorHandling ErrorHandling
	output        io.Writer
	parsed        bool
	args          []string
	flags         map[string]*Flag
	set           map[string]struct{}
}

// NewFlagSet returns a new flag set with the given name and error handling policy.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	return &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
}

// Name returns the name of this flag set as given to [NewFlagSet] or [FlagSet.Init].
func (f *FlagSet) Name() string {
	return f.name
}

// ErrorHandling returns the error handling policy for this flag set.
func (f *FlagSet) ErrorHandling() ErrorHandling {
	return f.errorHandling
}

// Output returns the [io.Writer] to which usage information is written, according to the [ErrorHandling] policy. The
// writer returned is the same given to [NewFlagSet] or [FlagSet.SetOutput].
func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}

	return f.output
}

// SetOutput sets the [io.Writer] to use when writing usage information, according to the [ErrorHandling] policy.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}

// Init sets the name and error handling policy for this flag set.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
}

// PrintDefaults prints usage information and default values for all flags of this flag set to the output location
// configured with [NewFlagSet] or [FlagSet.SetOutput].
func (f *FlagSet) PrintDefaults() {
	f.VisitAll(func(flg *Flag) {
		var err error

		name, usage := UnquoteUsage(flg)

		if len(flg.Name) == 1 {
			_, err = fmt.Fprintf(f.Output(), "   -%s", flg.Name)
		} else {
			_, err = fmt.Fprintf(f.Output(), "  --%s", flg.Name)
		}

		if err != nil {
			panic(err)
		}

		if len(name) > 0 {
			_, err = fmt.Fprintf(f.Output(), " <%s>", name)
		}

		if err != nil {
			panic(err)
		}

		if len(flg.DefValue) > 0 {
			if _, ok := flg.Value.(*stringT); ok {
				_, err = fmt.Fprintf(f.Output(), " (default %q)", flg.DefValue)
			} else {
				_, err = fmt.Fprintf(f.Output(), " (default %s)", flg.DefValue)
			}
		}

		if err != nil {
			panic(err)
		}

		_, err = fmt.Fprintf(f.Output(), "\n        %s\n", usage)

		if err != nil {
			panic(err)
		}
	})
}

// PrintDefaults calls [FlagSet.PrintDefaults] for the global [CommandLine].
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// Lookup returns a [Flag] with the given name, or nil if no such flag exists.
func (f *FlagSet) Lookup(name string) *Flag {
	if flg, ok := f.flags[name]; ok {
		return flg
	}

	return nil
}

// Arg returns the i'th remaining argument after calling [FlagSet.Parse]. Returns an empty string if the argument does
// not exist, or [FlagSet.Parse] was not called.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}

	return f.args[i]
}

// Arg returns [FlagSet.Arg] for the global [CommandLine].
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// NArg returns the number of non-flag arguments remaining after calling [FlagSet.Parse].
func (f *FlagSet) NArg() int {
	return len(f.args)
}

// NArg returns [FlagSet.NArg] for the global [CommandLine].
func NArg() int {
	return CommandLine.NArg()
}

// Args returns a slice of non-flag arguments remaining after calling [FlagSet.Parse].
func (f *FlagSet) Args() []string {
	return f.args
}

// Args returns [FlagSet.Args] for the global [CommandLine].
func Args() []string {
	return CommandLine.Args()
}

// NFlag returns the number of flags in this flag set that have been set.
func (f *FlagSet) NFlag() int {
	return len(f.set)
}

// NFlag returns [FlagSet.NFlag] for the global [CommandLine].
func NFlag() int {
	return CommandLine.NFlag()
}

// Parsed returns whether or not [FlagSet.Parse] has been invoked on this flag set.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// Parsed returns [FlagSet.Parsed] for the global [CommandLine].
func Parsed() bool {
	return CommandLine.Parsed()
}

// Set updates the value of a flag with the given string. Returns an error if the flag doesn't exist or the value is
// invalid.
func (f *FlagSet) Set(name, value string) error {
	flg := f.Lookup(name)
	if flg == nil {
		return fmt.Errorf("flag '%s' does not exist", name)
	}

	err := flg.Value.Set(value)
	if err != nil {
		return err
	}

	if f.set == nil {
		f.set = make(map[string]struct{})
	}

	f.set[name] = struct{}{}
	return nil
}

// Set returns [FlagSet.Set] for the global [CommandLine].
func Set(name, value string) error {
	return CommandLine.Set(name, value)
}

// Parse processes the given arguments and updates the flags of this flag set. The arguments given should not include
// the command name. Parse should only be called after all flags have been registered and before flags are accessed by
// the application.
//
// The return value will be [ErrHelp] if -help or -h were set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	usage := f.Usage
	if usage == nil {
		usage = f.PrintDefaults
	}

	err := f.parse(arguments)
	if err == nil {
		return nil
	}

	if f.errorHandling == ContinueOnError {
		usage()
		return err
	}

	if f.errorHandling == PanicOnError {
		usage()
		panic(err)
	}

	if errors.Is(err, ErrHelp) && f.errorHandling == ExitOnError {
		usage()
		os.Exit(0)
		return nil
	}

	usage()
	os.Exit(2)
	return nil
}

// Set invokes [FlagSet.Parse] for the global [CommandLine] with arguments from [os.Args].
func Parse() error {
	return CommandLine.Parse(os.Args[1:])
}

// VisitAll traverses all set flags in lexical order and executes fn for each one.
func (f *FlagSet) Visit(fn func(*Flag)) {
	for _, name := range slices.Sorted(maps.Keys(f.set)) {
		fn(f.Lookup(name))
	}
}

// Visit calls [FlagSet.Visit] for the global [CommandLine].
func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}

// VisitAll traverses all flags in lexical order and executes fn for each one.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, name := range slices.Sorted(maps.Keys(f.flags)) {
		fn(f.Lookup(name))
	}
}

// VisitAll calls [FlagSet.VisitAll] for the global [CommandLine].
func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

// Var registers a flag with an arbitrary [Value].
func (f *FlagSet) Var(value Value, name string, usage string) {
	if f.flags == nil {
		f.flags = make(map[string]*Flag)
	}

	if strings.HasPrefix(name, "-") {
		panic(fmt.Sprintf("flag '%s' has invalid name beginning with '-'", name))
	}
	if strings.HasSuffix(name, "-") {
		panic(fmt.Sprintf("flag '%s' has invalid name ending with '-'", name))
	}

	invalid := strings.ContainsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '.'
	})
	if invalid {
		panic(fmt.Sprintf("flag '%s' has invalid name (must be alphanumeric and any of '-.')", name))
	}

	if f.Lookup(name) != nil {
		panic(fmt.Sprintf("flag '%s' already registered with this flag set", name))
	}

	f.flags[name] = &Flag{
		Name:     name,
		Usage:    usage,
		Value:    value,
		DefValue: value.String(),
	}
}

// Set calls [FlagSet.Var] for the global [CommandLine].
func Var(value Value, name string, usage string) {
	CommandLine.Var(value, name, usage)
}

func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	var p bool
	f.Var(newBoolT(value, &p), name, usage)

	return &p
}

func Bool(name string, value bool, usage string) *bool {
	return CommandLine.Bool(name, value, usage)
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.Var(newBoolT(value, p), name, usage)
}

func BoolVar(p *bool, name string, value bool, usage string) {
	CommandLine.BoolVar(p, name, value, usage)
}

func (f *FlagSet) String(name string, value string, usage string) *string {
	var p string
	f.Var(newStringT(value, &p), name, usage)

	return &p
}

func String(name string, value string, usage string) *string {
	return CommandLine.String(name, value, usage)
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.Var(newStringT(value, p), name, usage)
}

func StringVar(p *string, name string, value string, usage string) {
	CommandLine.StringVar(p, name, value, usage)
}

func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	var p time.Duration
	f.Var(newDurationT(value, &p), name, usage)

	return &p
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return CommandLine.Duration(name, value, usage)
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.Var(newDurationT(value, p), name, usage)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	CommandLine.DurationVar(p, name, value, usage)
}

func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	var p float64
	f.Var(newFloat64T(value, &p), name, usage)

	return &p
}

func Float64(name string, value float64, usage string) *float64 {
	return CommandLine.Float64(name, value, usage)
}

func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	f.Var(newFloat64T(value, p), name, usage)
}

func Float64Var(p *float64, name string, value float64, usage string) {
	CommandLine.Float64Var(p, name, value, usage)
}

func (f *FlagSet) Int(name string, value int, usage string) *int {
	var p int
	f.Var(newIntT(value, &p), name, usage)

	return &p
}

func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.Var(newIntT(value, p), name, usage)
}

func IntVar(p *int, name string, value int, usage string) {
	CommandLine.IntVar(p, name, value, usage)
}

func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	var p int64
	f.Var(newInt64T(value, &p), name, usage)

	return &p
}

func Int64(name string, value int64, usage string) *int64 {
	return CommandLine.Int64(name, value, usage)
}

func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.Var(newInt64T(value, p), name, usage)
}

func Int64Var(p *int64, name string, value int64, usage string) {
	CommandLine.Int64Var(p, name, value, usage)
}

func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	var p uint
	f.Var(newUintT(value, &p), name, usage)

	return &p
}

func Uint(name string, value uint, usage string) *uint {
	return CommandLine.Uint(name, value, usage)
}

func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	f.Var(newUintT(value, p), name, usage)
}

func UintVar(p *uint, name string, value uint, usage string) {
	CommandLine.UintVar(p, name, value, usage)
}

func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	var p uint64
	f.Var(newUint64T(value, &p), name, usage)

	return &p
}

func Uint64(name string, value uint64, usage string) *uint64 {
	return CommandLine.Uint64(name, value, usage)
}

func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.Var(newUint64T(value, p), name, usage)
}

func Uint64Var(p *uint64, name string, value uint64, usage string) {
	CommandLine.Uint64Var(p, name, value, usage)
}

func (f *FlagSet) BoolFunc(name, usage string, fn func(string) error) {
	f.Var(boolFuncT(fn), name, usage)
}

func BoolFunc(name, usage string, fn func(string) error) {
	CommandLine.BoolFunc(name, usage, fn)
}

func (f *FlagSet) Func(name, usage string, fn func(string) error) {
	f.Var(funcT(fn), name, usage)
}

func Func(name, usage string, fn func(string) error) {
	CommandLine.Func(name, usage, fn)
}

func (f *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	f.Var(newTextT(value, p), name, usage)
}

func TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	CommandLine.TextVar(p, name, value, usage)
}

func (f *FlagSet) parse(arguments []string) error {
	var err error

	f.parsed = true

	for len(arguments) > 0 {
		arg := arguments[0]

		if arg == "-" {
			f.args = arguments[0:]
			return nil
		}
		if arg == "--" {
			f.args = arguments[1:]
			return nil
		}

		long, ok := strings.CutPrefix(arg, "--")
		if ok {
			arguments, err = f.parseLong(long, arguments[1:])
			if err != nil {
				return err
			}

			continue
		}

		short, ok := strings.CutPrefix(arg, "-")
		if ok {
			arguments, err = f.parseShort(short, arguments[1:])
			if err != nil {
				return err
			}

			continue
		}

		f.args = arguments
		return nil
	}

	f.args = arguments
	return nil
}

func (f *FlagSet) parseLong(arg string, arguments []string) ([]string, error) {
	arg, value, inlineVal := strings.Cut(arg, "=")

	flg := f.Lookup(arg)
	if flg == nil && arg == "help" {
		return nil, ErrHelp
	}
	if flg == nil {
		return nil, fmt.Errorf("flag '--%s' does not exist", arg)
	}

	if isBoolFlag(flg) {
		if !inlineVal {
			value = "true"
		}
	} else {
		if !inlineVal {
			if len(arguments) == 0 {
				return nil, fmt.Errorf("missing argument to flag '--%s'", arg)
			}

			arguments, value = arguments[1:], arguments[0]
		}
	}

	if err := f.Set(arg, value); err != nil {
		return nil, err
	}

	return arguments, nil
}

func (f *FlagSet) parseShort(short string, arguments []string) ([]string, error) {
	for len(short) > 0 {
		args := strings.SplitN(short, "", 2)

		if len(args) == 1 {
			short = ""
		}
		if len(args) == 2 {
			short = args[1]
		}

		flg := f.Lookup(args[0])
		if flg == nil && args[0] == "h" {
			return nil, ErrHelp
		}
		if flg == nil {
			return nil, fmt.Errorf("flag '-%s' does not exist", args[0])
		}

		if isBoolFlag(flg) {
			if err := f.Set(args[0], "true"); err != nil {
				return nil, err
			}
		} else {
			if short != "" {
				// rest is arg
				if err := f.Set(args[0], short); err != nil {
					return nil, err
				}
			} else {
				// take next arg
				if len(arguments) == 0 {
					return nil, fmt.Errorf("missing argument to flag '-%s'", args[0])
				}

				if err := f.Set(args[0], arguments[0]); err != nil {
					return nil, err
				}

				arguments = arguments[1:]
			}

			return arguments, nil
		}
	}

	return arguments, nil
}
