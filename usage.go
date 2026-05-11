package cmder

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"text/template"

	"github.com/brandon1024/cmder/getopt"
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

	{{- print (flag_usage .) -}}
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
	tmpl, err := template.New("usage").Funcs(funcs(ops)).Parse(ops.usageTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(ops.outputWriter, cmd)
}

// help renders extended help text for a [Command].
func help(cmd command, ops *ExecuteOptions) error {
	tmpl, err := template.New("help").Funcs(funcs(ops)).Parse(ops.helpTemplate)
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
//   - flags(c):               Return the flagset of c.
//   - flag_usage(fs):         Return the rendered flag usage for the given flagset.
//   - lower(str):             Return string argument in lowercase.
//   - upper(str):             Return string argument in uppercase.
//   - split(str):             Split a string.
//   - replace(str, old, new): Replace occurrences of a string.
//   - join(slice, delim):     Join a list of strings.
//   - contains(str, other):   Check if a string contains another string
//   - trim(str):              Trim all leading and trailing whitespace of str.
//   - lines(str):             Split str into a slice of text lines.
func funcs(ops *ExecuteOptions) template.FuncMap {
	return template.FuncMap{
		"commands":   subcommands,
		"flags":      flags(ops),
		"flag_usage": flagUsage,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"split":      strings.Split,
		"replace":    strings.ReplaceAll,
		"join":       strings.Join,
		"contains":   strings.Contains,
		"trim":       strings.TrimSpace,
		"lines":      strings.Lines,
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

// flags returns a template func which produces a flagset (either a standard [flag.FlagSet] or [getopt.PosixFlagSet])
// according to the options defines in ops.
func flags(ops *ExecuteOptions) func(cmd command) any {
	return func(cmd command) any {
		if ops.nativeFlags {
			return cmd.fs
		}

		return &getopt.PosixFlagSet{FlagSet: cmd.fs, RelaxedParsing: ops.relaxedFlags}
	}
}

// flagsetPrinter is a flagset (either [flag.FlagSet] or [getopt.PosixFlagSet]) which can render its usage.
type flagsetPrinter interface {
	PrintDefaults()
	Output() io.Writer
	SetOutput(io.Writer)
}

// flagUsage returns the text rendered by either [flag.FlagSet.PrintDefaults] or [getopt.PosixFlagSet.PrintDefaults].
func flagUsage(fs flagsetPrinter) string {
	var buf bytes.Buffer

	original := fs.Output()
	defer fs.SetOutput(original)

	fs.SetOutput(&buf)
	fs.PrintDefaults()

	return buf.String()
}
