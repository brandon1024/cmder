package cmder

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
)

// The default text template for rendering command usage information.
const DefaultUsageTemplate = CobraUsageTemplate

// Text template for rendering command usage information in a format similar to that of the popular
// [github.com/spf13/cobra] library.
const CobraUsageTemplate = `
{{- define "HelpText" -}}
{{- trim .Cmd.HelpText -}}
{{- end -}}

{{- define "UsageText" }}

Usage:
{{ shift (trim .Cmd.UsageLine) 2 }}
{{- end -}}

{{- define "ExampleText" -}}
{{- if .Cmd.ExampleText }}

Examples:
{{ shift (trim .Cmd.ExampleText) 2 }}
{{- end -}}
{{- end -}}

{{- define "AvailableCommandsText" -}}
{{- with .SubcommandSummary }}

Available Commands:
{{ shift (table . 2) 2 }}
{{- end -}}
{{- end -}}

{{- define "FlagsText" -}}
{{- with .FlagSummary }}

Flags:
{{ shift (table . 3) 2 }}
{{- end -}}
{{- end -}}

{{- define "AdditionalCommandInfoText" -}}
{{- if .SubcommandSummary }}

Use "{{ .Cmd.Name }} [command] --help" for more information about a command.
{{- end -}}
{{- end -}}


{{- template "HelpText" . }}
{{- template "UsageText" . }}
{{- template "ExampleText" . }}
{{- template "AvailableCommandsText" . }}
{{- template "FlagsText" . }}
{{- template "AdditionalCommandInfoText" . }}
`

// Text template for rendering command usage information in a minimal format similar to that of the [flag] standard
// library.
const StdFlagUsageTemplate = `
TODO
`

// Text template for rendering command usage information in an extended format similar to the format used by kubectl.
// This format is suitable for commands that provide more thorough documentation (especially for flags).
const LongUsageTemplate = `
TOOD
`

// The default writer for command usage information. Standard error is recommended, but you can override this if needed
// (particularly useful in tests).
var UsageOutputWriter io.Writer = os.Stderr

// Data passed to the usage text template.
type usageTemplateData struct {
	Cmd               Command
	SubcommandSummary [][]string
	FlagSummary       [][]string
}

// A group of one or two flags, typically a shorthand and a longer flag (e.g. -h/--help).
type flagGroup struct {
	// The short flag names.
	Short []string

	// The long flag names.
	Long []string

	// A descriptive representation of the flag(s).
	Description string

	// The default value for the flag(s).
	DefValue string
}

// Render usage text for the given [Command] using the default template [DefaultUsageTemplate]. Output is written to
// [UsageOutputWriter].
//
// The following template functions are available:
//
//	table([][]string): render a table
//	shift(string, int): left pad a string by the given number of spaces
//	trim(string): trim leading and trailing whitespace from the string
func RenderUsage(cmd Command, fs *flag.FlagSet) error {
	data := usageTemplateData{
		Cmd:               cmd,
		SubcommandSummary: [][]string{},
		FlagSummary:       buildFlagSummary(fs),
	}

	if c, ok := cmd.(RootCommand); ok {
		for _, sub := range c.Subcommands() {
			if sub.Hidden() {
				continue
			}

			data.SubcommandSummary = append(data.SubcommandSummary, []string{
				sub.Name(),
				sub.ShortHelpText(),
			})
		}
	}

	tmpl, err := template.New(cmd.Name()).Funcs(map[string]any{
		"table": tableTemplateFunc,
		"shift": shiftTemplateFunc,
		"trim":  strings.TrimSpace,
	}).Parse(DefaultUsageTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(UsageOutputWriter, data)
}

// Build a representation of all flags in the given [*flag.FlagSet].
//
// Command-line tools often provide abbreviated single-character flags in addition to longer flag names. Common examples
// are `-h` and `--help`, or `-n` and `--count`. The standard [flag] package makes no differentiation between short and
// long flags -- they are all distinct flags. In order to group short/long flags in usage output, this helper iterates
// over the flag set and builds a set of flag groups.
//
// Each flag group can have one or more names. The names in the flag group are sorted in lexicographical order.
//
// A flag group is comprised of one or two flags that:
//
// - have the same usage text ([flag.Usage])
// - have the same default value ([flag.Default])
func buildFlagSummary(fs *flag.FlagSet) [][]string {
	// groups of flags, keyed by "f.Usage+f.DefaultValue"
	groups := map[string][]*flag.Flag{}

	// first, collect flags into groups
	fs.VisitAll(func(f *flag.Flag) {
		key := f.Usage + f.DefValue

		group, ok := groups[key]
		if !ok {
			group = []*flag.Flag{}
		}

		groups[key] = append(group, f)
	})

	// next, associate short flags with long flags
	result := [][]string{}

	for _, flags := range groups {
		// divide the flags into two groups: short and long flags
		short := []*flag.Flag{}
		long := []*flag.Flag{}

		for _, f := range flags {
			if len(f.Name) > 1 {
				long = append(long, f)
			} else {
				short = append(short, f)
			}
		}

		// one by one, match up the flags
		for i := range max(len(short), len(long)) {
			switch {
			case i < len(short) && i < len(long):
				result = append(result, []string{
					fmt.Sprintf("-%s, --%s", short[i].Name, long[i].Name),
					fmt.Sprintf("%s (default %s)", short[i].Usage, short[i].DefValue),
				})
			case i < len(short):
				result = append(result, []string{
					fmt.Sprintf("-%s", short[i].Name),
					fmt.Sprintf("%s (default %s)", short[i].Usage, short[i].DefValue),
				})
			case i < len(long):
				result = append(result, []string{
					fmt.Sprintf("--%s", long[i].Name),
					fmt.Sprintf("%s (default %s)", long[i].Usage, long[i].DefValue),
				})
			}
		}
	}

	return result
}

// Render the given table of data in a nicely formatted column format.
func tableTemplateFunc(data [][]string, padding int) string {
	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, padding, ' ', uint(0))
	for _, cols := range data {
		fmt.Fprintf(w, "%s\n", strings.Join(cols, "\t"))
	}

	w.Flush()

	return strings.TrimSpace(buf.String())
}

// Shift each line in data by the given number of spaces.
func shiftTemplateFunc(data string, spaces int) string {
	var lines []string
	for line := range strings.Lines(data) {
		lines = append(lines, strings.Repeat(" ", spaces)+line)
	}

	return strings.Join(lines, "")
}
