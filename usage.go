package cmder

import (
	"bytes"
	"cmp"
	"flag"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/brandon1024/cmder/getopt"
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

		{{ with (index . 0).DefValue }}
			{{- printf " (default %s)" . -}}
		{{- end -}}

		{{- println -}}

		{{- printf "      %s\n" (index (unquote (index . 0)) 1) -}}
	{{- end -}}
{{- end -}}

{{- if (commands .Command) -}}
	{{- println -}}
	{{- printf "Use \"%s [command] --help\" for more information about a command.\n" .Command.Name -}}
{{- end -}}`

// Text template for rendering command usage information in a minimal format similar to that of the [flag] library.
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
		"unquote":   unquote,
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

// flagUsage dumps the flag usage as rendered by the flag library. See [flag.FlagSet.PrintDefaults].
func flagUsage(cmd command) string {
	out := cmd.fs.Output()
	defer cmd.fs.SetOutput(out)

	var buf bytes.Buffer
	cmd.fs.SetOutput(&buf)

	cmd.fs.PrintDefaults()

	return buf.String()
}

// unquote calls [flag.UnquoteUsage] for the given [flag.Flag].
func unquote(flg *flag.Flag) []string {
	if isBoolFlag(flg) {
		return []string{"", flg.Usage}
	}

	name, usage := flag.UnquoteUsage(flg)

	// if no backquoted names found, try to infer from [flag.Getter]
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

// isHiddenFlag checks if the given flag is hidden.
func isHiddenFlag(flg *flag.Flag) bool {
	hf, ok := flg.Value.(getopt.HiddenFlag)
	return ok && hf.IsHiddenFlag()
}

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// isBoolFlag checks if the given flag is a boolean flag.
func isBoolFlag(flg *flag.Flag) bool {
	hf, ok := flg.Value.(boolFlag)
	return ok && hf.IsBoolFlag()
}
