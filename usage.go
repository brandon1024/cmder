package cmder

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"text/template"
)

// Text template for rendering command usage information in a format similar to that of the popular
// [github.com/spf13/cobra] library.
const CobraUsageTemplate = `{{ trim .Command.HelpText }}

Usage:
  {{ trim .Command.UsageLine }}

Examples:
{{ range (lines (trim .Command.ExampleText)) }}  {{ . }}{{ end }}
{{- println -}}

{{- with (commands .Command) -}}
	{{- println -}}
	{{- println "Available Commands:" -}}
	{{- range . -}}
		{{- printf "  %-13s  %s\n" .Name .ShortHelpText -}}
	{{- end -}}
{{- end -}}

{{- with (flags .) -}}
	{{- println -}}
	{{- println "Flags:" -}}
	{{- range . -}}
		{{- printf "  %-32s  %s%s\n" (printf "%s%s" (or (and (eq (len .Name) 1) " -") "--") .Name) .Usage (and .DefValue (printf " (default \"%s\")" .DefValue)) -}}
	{{- end -}}
{{- end -}}

{{- if (commands .Command) -}}
	{{- println -}}
	{{- printf "Use \"%s [command] --help\" for more information about a command.\n" .Command.Name -}}
{{- end -}}`

// Text template for rendering command usage information in a minimal format similar to that of the [flag] standard
// library.
const StdFlagUsageTemplate = `usage: {{ .Command.UsageLine }}
{{ flagusage . }}`

// The text template for rendering command usage information.
var UsageTemplate = CobraUsageTemplate

// The default writer for command usage information. Standard error is recommended, but you can override this if needed
// (particularly useful in tests).
var UsageOutputWriter io.Writer = os.Stderr

// usage renders usage text for a [Command] using the default template [UsageTemplate]. Output is written to
// [UsageOutputWriter].
func usage(cmd command) error {
	tmpl, err := template.New("usage").Funcs(template.FuncMap{
		"commands":  collectSubcommands,
		"flags":     flags,
		"flagusage": flagUsage,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"split":     strings.Split,
		"replace":   strings.ReplaceAll,
		"join":      strings.Join,
		"contains":  strings.Contains,
		"trim":      strings.TrimSpace,
		"lines":     strings.Lines,
	}).Parse(UsageTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(UsageOutputWriter, cmd)
}

// collect a slice of flags of the given command.
func flags(cmd command) []*flag.Flag {
	var collected []*flag.Flag

	cmd.fs.VisitAll(func(f *flag.Flag) {
		collected = append(collected, f)
	})

	return collected
}

// flagUsage dumps the flag usage as rendered by the stdlib.
func flagUsage(cmd command) string {
	out := cmd.fs.Output()
	defer cmd.fs.SetOutput(out)

	var buf bytes.Buffer
	cmd.fs.SetOutput(&buf)

	cmd.fs.PrintDefaults()

	return buf.String()
}
