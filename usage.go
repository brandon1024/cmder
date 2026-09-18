package cmder

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/brandon1024/cmder/getopt"
)

// DefaultHelpTemplate is a text template for rendering extended command help information.
//
// This help template is organized into sections, and each section can be overridden with user-provided content. See
// [WithNamedTemplate]. The following sections can be overridden:
//
//   - section.help
//   - section.header
//   - section.synopsis
//   - section.examples
//   - section.subcommands
//   - section.options
//   - section.see-also
//   - section.footer
//
// This template is used to render help text by default. You may configure an alternate template with
// [WithHelpTemplate].
const DefaultHelpTemplate = `
{{- block "section.help" . -}}
	{{- println (trim .Command.HelpText) -}}
	{{- println -}}
{{- end -}}` + DefaultUsageTemplate

// DefaultUsageTemplate is a text template for rendering command usage information.
//
// This usage template is organized into sections, and each section can be overridden with user-provided content. See
// [WithNamedTemplate]. The following sections can be overridden:
//
//   - section.header
//   - section.synopsis
//   - section.examples
//   - section.subcommands
//   - section.options
//   - section.see-also
//   - section.footer
//
// This template is used to render usage text by default. You may configure an alternate template with
// [WithUsageTemplate].
const DefaultUsageTemplate = `
{{- block "section.header" . -}}{{- end -}}

{{- block "section.synopsis" . -}}
	{{- println "Usage:" -}}
	{{- printf "  %s" (trim .Command.UsageLine) -}}
	{{- println -}}
{{- end -}}

{{- block "section.examples" . -}}
	{{- with .Command.ExampleText -}}
		{{- println -}}
		{{- println "Examples:" -}}
		{{- range (lines (trim .)) -}}
			{{- printf "  %s" . -}}
		{{- end -}}
		{{- println -}}
	{{- end -}}
{{- end -}}

{{- block "section.subcommands" . -}}
	{{- with (commands .) -}}
		{{- println -}}
		{{- println "Available Commands:" -}}
		{{- range . -}}
			{{- printf "  %-13s  %s\n" .Name .ShortHelpText -}}
		{{- end -}}
	{{- end -}}
{{- end -}}

{{- block "section.options" . -}}
	{{- with (flags .) -}}
		{{- println -}}
		{{- println "Flags:" -}}

		{{- print (flag_usage .) -}}
	{{- end -}}

	{{- with (parents .) -}}
		{{- range . -}}
			{{- if (flags .) -}}
				{{- println -}}
				{{- printf "Flags for \"%s\":\n" .Command.Name -}}

				{{- print (flag_usage (flags .)) -}}
			{{- end -}}
		{{- end -}}
	{{- end -}}
{{- end -}}

{{- block "section.see-also" . -}}
	{{- if (commands .) -}}
		{{- println -}}
		{{- printf "Use \"%s [command] --help\" for more information about a command.\n" .Command.Name -}}
	{{- end -}}
{{- end -}}

{{- block "section.footer" . -}}{{- end -}}
`

// ErrShowUsage instructs cmder to render usage.
var ErrShowUsage = errors.New("cmder: usage requested")

// ErrShowHelp instructs cmder to render help.
var ErrShowHelp = errors.New("cmder: help requested")

// usage renders usage text for a [Command].
func usage(cmd command, ops *ExecuteOptions) error {
	return render(cmd, ops, ops.usageTemplate)
}

// help renders extended help text for a [Command].
func help(cmd command, ops *ExecuteOptions) error {
	return render(cmd, ops, ops.helpTemplate)
}

// render the given template string for cmd.
func render(cmd command, ops *ExecuteOptions, renderTmpl string) error {
	tmpl, err := template.New("_out").Funcs(funcs(ops)).Parse(renderTmpl)
	if err != nil {
		return err
	}

	for name, def := range ops.secondaryTemplates {
		_, err = tmpl.New(name).Parse(def)
		if err != nil {
			return err
		}
	}

	return tmpl.ExecuteTemplate(ops.outputWriter, "_out", cmd)
}

// funcs returns template functions which can be used in usage/help text templates.
//
// The following template functions are available:
//
//   - commands(c):            Collect all subcommands of c into a map, keyed by name, which are not hidden.
//   - parents(c):             Return a slice of all parent commands of c, in descending order of depth.
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
		"parents":    parents,
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
// according to the options defines in ops. Returns nil if the command doesn't implement [FlagInitializer].
func flags(ops *ExecuteOptions) func(cmd command) any {
	return func(cmd command) any {
		if _, ok := cmd.Command.(FlagInitializer); !ok {
			return nil
		}

		if ops.nativeFlags {
			return cmd.fs
		}

		return &getopt.PosixFlagSet{FlagSet: cmd.fs, RelaxedParsing: ops.relaxedFlags}
	}
}

// parents returns a slice of all parent commands of cmd. If the command is a root command, returns an empty slice.
func parents(cmd command) []command {
	parents := slices.Clone(cmd.parents)
	slices.Reverse(parents)

	return parents
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
