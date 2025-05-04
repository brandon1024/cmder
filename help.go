package cmder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
)

// The default text template for rendering help/usage information.
const DefaultHelpTemplate = CobraHelpTemplate

// Text template for rendering help/usage information in a format similar to that of the popular
// [github.com/spf13/cobra] library.
const CobraHelpTemplate = `
{{- define "HelpText" -}}
{{- trim .HelpText -}}
{{- end -}}

{{- define "UsageText" }}

Usage:
{{ shift (trim .UsageLine) 2 }}
{{- end -}}

{{- define "ExampleText" -}}
{{- if .ExampleText }}

Examples:
{{ shift (trim .ExampleText) 2 }}
{{- end -}}
{{- end -}}

{{- define "AvailableCommandsText" -}}
{{- with .Subcommands }}

Available Commands:
{{ shift (table .) 2 }}
{{- end -}}
{{- end -}}

{{- define "FlagsText" -}}
{{- with .Flags }}

Flags:
{{ shift (table .) 2 }}
{{- end -}}
{{- end -}}

{{- define "AdditionalCommandInfoText" -}}
{{- if .Subcommands }}

Use "{{ .Name }} [command] --help" for more information about a command.
{{- end -}}
{{- end -}}


{{- template "HelpText" . }}
{{- template "UsageText" . }}
{{- template "ExampleText" . }}
{{- template "AvailableCommandsText" . }}
{{- template "FlagsText" . }}
{{- template "AdditionalCommandInfoText" . }}
`

// Text template for rendering help/usage information in a minimal format similar to that of the [flag] standard library.
const StdFlagHelpTemplate = `
TODO
`

// Text template for rendering help/usage information in an extended format similar to the format used by kubectl. This
// format is suitable for commands that provide more thorough documentation (especially for flags).
const LongHelpTemplate = `
TOOD
`

// The default writer for help/usage information. Standard error is recommended, but you can override this if needed
// (particularly useful in tests).
var HelpOutputWriter io.Writer = os.Stderr

type helpTemplateData struct {
	Name        string
	UsageLine   string
	HelpText    string
	ExampleText string
	Subcommands [][]string
	Flags       [][]string
}

// Render help text for the given [Command] using the default template [DefaultHelpTemplate]. Output is written to
// [HelpOutputWriter].
//
// The following template functions are available:
//
//	table([][]string): render a table
//	shift(string, int): left pad a string by the given number of spaces
//	trim(string): trim leading and trailing whitespace from the string
func RenderHelp(cmd Command) error {
	data := helpTemplateData{
		Name:        cmd.Name(),
		UsageLine:   cmd.UsageLine(),
		HelpText:    cmd.HelpText(),
		ExampleText: cmd.ExampleText(),
		Subcommands: [][]string{},
		Flags:       [][]string{},
	}

	if c, ok := cmd.(RootCommand); ok {
		for _, sub := range c.Subcommands() {
			if sub.Hidden() {
				continue
			}

			data.Subcommands = append(data.Subcommands, []string{
				sub.Name(),
				sub.ShortHelpText(),
			})
		}
	}

	tmpl, err := template.New(cmd.Name()).Funcs(map[string]any{
		"table": tableTemplateFunc,
		"shift": shiftTemplateFunc,
		"trim":  strings.TrimSpace,
	}).Parse(DefaultHelpTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(HelpOutputWriter, data)
}

func tableTemplateFunc(data [][]string) string {
	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', uint(0))
	for _, line := range data {
		fmt.Fprintf(w, "%s\n", strings.Join(line, "\t"))
	}

	w.Flush()

	return strings.TrimSpace(buf.String())
}

func shiftTemplateFunc(data string, spaces int) string {
	var lines []string
	for line := range strings.Lines(data) {
		lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", spaces), line))
	}

	return strings.Join(lines, "")
}
