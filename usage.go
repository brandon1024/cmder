package cmder

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"text/template"
	"time"
)

// DefaultHelpTemplate is a text template for rendering extended command help information.
const DefaultHelpTemplate = `{{ trim .Command.HelpText }}{{ println }}{{ println }}` + DefaultUsageTemplate

// DefaultUsageTemplate is a text template for rendering command usage information.
const DefaultUsageTemplate = `Usage:
{{- println -}}
{{- printf "  %s" (trim .Command.UsageLine) -}}
{{- println -}}

{{- with .Command.ExampleText -}}
	{{- println -}}
	{{- println "Examples:" -}}
	{{- range (lines (trim .)) -}}
		{{- printf "  %s" . -}}
	{{- end -}}
	{{- println -}}
{{- end -}}

{{- with (commands .) -}}
	{{- println -}}
	{{- println "Available Commands:" -}}
	{{- range . -}}
		{{- printf "  %-13s  %s\n" .Name .ShortHelpText -}}
	{{- end -}}
{{- end -}}

{{- with (flags .) -}}
	{{- println -}}
	{{- println "Flags:" -}}

	{{- $print_started := false -}}

	{{- range . -}}
		{{- if $print_started -}}
			{{- println -}}
		{{- end -}}
		{{- $print_started = true -}}

		{{- printf "  " -}}

		{{- range $index, $flg := . -}}
			{{- if (ne $index 0) -}}
				{{- printf ", " -}}
			{{- end -}}

			{{- if (eq (len $flg.Name) 1) -}}
				{{- printf "-%s" .Name -}}
			{{- else -}}
				{{- printf "--%s" .Name -}}
			{{- end -}}

			{{- $name := (index (unquote $flg) 0) -}}

			{{- if (and $name (eq (len $flg.Name) 1)) -}}
				{{- printf " <%s>" $name -}}
			{{- else if $name -}}
				{{- printf "=<%s>" $name -}}
			{{- end -}}
		{{- end -}}

		{{ if (not (zero (index . 0))) }}
			{{- printf " (default %s)" (index . 0).DefValue -}}
		{{- end -}}

		{{- println -}}

		{{- printf "      %s\n" (index (unquote (index . 0)) 1) -}}
	{{- end -}}
{{- end -}}

{{- if (commands .) -}}
	{{- println -}}
	{{- printf "Use \"%s [command] --help\" for more information about a command.\n" .Command.Name -}}
{{- end -}}`

// ErrShowUsage instructs cmder to render usage.
var ErrShowUsage = errors.New("cmder: usage requested")

// ErrShowHelp instructs cmder to render help.
var ErrShowHelp = errors.New("cmder: help requested")

// usage renders usage text for a [Command].
func usage(cmd command, ops *ExecuteOptions) error {
	tmpl, err := template.New("usage").Funcs(funcs()).Parse(ops.usageTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(ops.outputWriter, cmd)
}

// help renders extended help text for a [Command].
func help(cmd command, ops *ExecuteOptions) error {
	tmpl, err := template.New("help").Funcs(funcs()).Parse(ops.helpTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(ops.outputWriter, cmd)
}

// funcs returns template functions which can be used in usage/help text templates.
//
// The following template functions are available:
//
//   - commands(c):            Collect all subcommands of c into a map, keyed by name.
//   - flags(c):               Collect all flags of c, organized by flag group name.
//   - unquote(f):             Call UnquoteUsage on flag f.
//   - zero(f):                Check if the default value for f is the zero value or not.
//   - lower(str):             Return string argument in lowercase.
//   - upper(str):             Return string argument in uppercase.
//   - split(str):             Split a string.
//   - replace(str, old, new): Replace occurrences of a string.
//   - join(slice, delim):     Join a list of strings.
//   - contains(str, other):   Check if a string contains another string
//   - trim(str):              Trim all leading and trailing whitespace of str.
//   - lines(str):             Split str into a slice of text lines.
func funcs() template.FuncMap {
	return template.FuncMap{
		"commands": subcommands,
		"flags":    flags,
		"unquote":  unquote,
		"zero":     zero,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"split":    strings.Split,
		"replace":  strings.ReplaceAll,
		"join":     strings.Join,
		"contains": strings.Contains,
		"trim":     strings.TrimSpace,
		"lines":    strings.Lines,
	}
}

// subcommands returns a map of (visible) child subcommands for cmd.
func subcommands(cmd command) map[string]Command {
	subcommands := map[string]Command{}

	for name, c := range collectSubcommands(cmd.Command) {
		if hidden, ok := c.(HiddenCommand); !ok || !hidden.Hidden() {
			subcommands[name] = c
		}
	}

	return subcommands
}

// flags organizes the flags of cmd and returns them.
//
// The flags of cmd are grouped by [flag.Value] equivalence. This allows flags to be grouped together in the rendered
// usage text when two flags are aliases of each other. This is often the case for short flags which are aliases of
// longer flags (e.g. '-a' is an alias of '--all').
//
//	-a <string>, --addr=<string>
//	-s <string>, --serial-number=<string>
//
// The resulting map entries are keyed by the flag group name, which is the longest flag name in the group. The map
// values are slices of (one or more) flags in the flag group, sorted by flag name length ('-a' before '--all').
func flags(cmd command) map[string][]*flag.Flag {
	var collected []*flag.Flag

	cmd.fs.VisitAll(func(f *flag.Flag) {
		if !isHiddenFlag(f) {
			collected = append(collected, f)
		}
	})

	// sort flags by name length in descending order to ensure that keys in resulting map will use long names first
	slices.SortFunc(collected, func(a, b *flag.Flag) int {
		return cmp.Compare(len(b.Name), len(a.Name))
	})

	groups := map[string][]*flag.Flag{}

	for len(collected) > 0 {
		var flg *flag.Flag

		// pop the head of the slice
		flg, collected = collected[0], collected[1:]

		// update groups
		groups[flg.Name] = []*flag.Flag{flg}

		// traverse the flags again and find (and remove) any which match flg
		for i := len(collected) - 1; i >= 0; i-- {
			other := collected[i]

			if areSame(flg.Value, other.Value) {
				groups[flg.Name] = append(groups[flg.Name], other)
				collected = append(collected[:i], collected[i+1:]...)
			}
		}

		// sort by length (then lexical order), this time ascending (-a before --all)
		slices.SortFunc(groups[flg.Name], func(a, b *flag.Flag) int {
			if c := cmp.Compare(len(a.Name), len(b.Name)); c != 0 {
				return c
			}

			return cmp.Compare(a.Name, b.Name)
		})
	}

	return groups
}

// unquote calls [flag.UnquoteUsage] for the given [flag.Flag].
func unquote(flg *flag.Flag) []string {
	if isBoolFlag(flg) {
		return []string{"", flg.Usage}
	}

	name, usage := flag.UnquoteUsage(flg)

	// if no `names` found, try to infer from [flag.Getter]
	if name == "" {
		if g, ok := flg.Value.(flag.Getter); ok {
			switch g.Get().(type) {
			case uint, uint64:
				name = "uint"
			case int, int64:
				name = "int"
			case float64:
				name = "float"
			case time.Duration:
				name = "duration"
			default:
				name = "arg"
			}
		}
	}

	return []string{name, usage}
}

// zero checks if the default value of flg is the zero value for its type. This is used when rendering usage text
// to render default flag values only when the default value is interesting.
//
// This function expects that flg adhere's to the same requirements of the stdlib [flag] package, notably:
//
//	The flag package may call the String method with a zero-valued receiver, such as a nil pointer.
//
// Flags that don't respect this requirement will result in an error.
func zero(flg *flag.Flag) (ok bool, err error) {
	var z reflect.Value

	if typ := reflect.TypeOf(flg.Value); typ.Kind() == reflect.Pointer {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("cmder: flag '%s' is backed by a type that does not accept calling String() on the zero value (bug): %v",
				flg.Name, e)
		}
	}()

	return flg.DefValue == z.Interface().(flag.Value).String(), nil
}
