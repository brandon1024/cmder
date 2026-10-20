package internal

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// string funcs
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"split":    strings.Split,
		"replace":  strings.ReplaceAll,
		"join":     strings.Join,
		"contains": strings.Contains,
		"trim":     strings.TrimSpace,
		"lines":    strings.Lines,
		"wrap":     wrapLines,

		// flag funcs
		"flag_unquote":  unquote,
		"flag_zero":     zero,
		"flagset_usage": flagsetUsage,
	}
}

// wrapLines breaks s into lines, trying to maintain consistent line lengths under length characters.
func wrapLines(s string, length int) []string {
	var lines []string

	for line := range strings.Lines(strings.TrimSpace(s)) {
		line = strings.TrimSpace(line)

		for {
			breakpoint := len(line)

			for utf8.RuneCountInString(line) > length {
				space := strings.LastIndexFunc(line[:breakpoint], unicode.IsSpace)
				if space == -1 {
					break
				}

				breakpoint = space

				if utf8.RuneCountInString(line[:breakpoint]) <= length {
					break
				}
			}

			lines = append(lines, strings.TrimSpace(line[:breakpoint]))
			line = strings.TrimSpace(line[breakpoint:])

			if len(line) == 0 {
				break
			}
		}
	}

	return lines
}

// unquote is a wrapper over the standard [flag.UnquoteUsage] which returns a slice, allowing it to be used as a
// template func.
func unquote(flg *flag.Flag) []string {
	name, usage := flag.UnquoteUsage(flg)
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

	ok = flg.DefValue == z.Interface().(flag.Value).String()
	return
}

// flagsetPrinter is a flagset (either [flag.FlagSet] or [getopt.PosixFlagSet]) which can render its usage.
type flagsetPrinter interface {
	PrintDefaults()
	Output() io.Writer
	SetOutput(io.Writer)
}

// flagUsage returns the text rendered by either [flag.FlagSet.PrintDefaults] or [getopt.PosixFlagSet.PrintDefaults].
func flagsetUsage(fs flagsetPrinter) string {
	var buf bytes.Buffer

	original := fs.Output()
	defer fs.SetOutput(original)

	fs.SetOutput(&buf)
	fs.PrintDefaults()

	return buf.String()
}
