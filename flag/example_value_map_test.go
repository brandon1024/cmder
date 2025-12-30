package flag_test

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/brandon1024/cmder/flag"
)

// MapVar is a [flag.Value] for flags that accept map values. MapVar also implements [Flag.Getter].
//
// MapVar parses flag values which are key=value pairs. Multiple key=value pairs may be comma separated (e.g.
// key1=value1,key2=vlue2). Keys must be alphanumeric. Values can contain periods and whitespace if enclosed in double
// quotes. Unquoted whitespace will be trimmed.
//
//	key1=value1
//	key1=value1,key2=value2
//	key1 = "value, 1"
//	key1=v=1,key2=v=2
type MapVar map[string]string

// String returns the map, formatted as a set of key-value pairs:
func (m MapVar) String() string {
	var entries []string

	for _, k := range slices.Sorted(maps.Keys(m)) {
		entries = append(entries, k+"="+strconv.Quote(m[k]))
	}

	return strings.Join(entries, ",")
}

// Set fulfills the [flag.Value] interface. The given value must be a set of key-value pairs.
func (m MapVar) Set(value string) error {
	var (
		entries         = make(map[string]string)
		quoted, inValue bool
		key, val        string
	)

	for pos, c := range value {
		switch {
		case c == '"' && !inValue:
			return fmt.Errorf("illegal mapvar value at position %d (quoted key): %s", pos, value)
		case !quoted && c == '"':
			quoted = true
		case quoted && c == '"':
			quoted = false
		case quoted && inValue:
			val += string(c)
		case quoted && !inValue:
			return fmt.Errorf("illegal mapvar value at position %d (quoted key): %s", pos, value)
		case unicode.IsSpace(c):
		case !inValue && c == ',':
			return fmt.Errorf("illegal mapvar value at position %d (malformed pair missing value): %s", pos, value)
		case inValue && c == '=':
			val += string(c)
		case c == '=':
			inValue = true
		case c == ',':
			entries[key] = val
			key, val, inValue = "", "", false
		case inValue:
			val += string(c)
		case !inValue:
			key += string(c)
		default:
			panic(fmt.Errorf("bug: unhandled case: pos=%d, c=%v, quoted=%v, inValue=%v, key='%s', val='%s'",
				pos, c, quoted, inValue, key, val))
		}
	}

	if quoted {
		return fmt.Errorf("illegal mapvar value (quote mismatch): %s", value)
	}
	if !inValue {
		return fmt.Errorf("illegal mapvar value (malformed pair missing value): %s", value)
	}

	entries[key] = val

	for k, v := range entries {
		m[k] = v
	}

	return nil
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a
// map[string]string.
func (m MapVar) Get() any {
	return map[string]string(m)
}

// This example implements a custom flag type for string maps. You'll often find map flags on commands that perform
// templating of text files, for example.
func ExampleValue_map() {
	variables := MapVar{}

	fs := flag.NewFlagSet("map", flag.ContinueOnError)
	fs.Var(variables, "variable", "specify runtime variables")
	fs.Var(variables, "v", "specify runtime variables")

	args := []string{
		"--variable", "key1=value1",
		"-v", "key2=value2,key3=value3",
		`--variable=hello=" HI, WORLD "`,
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("variables: %s\n", variables.String())
	// Output:
	// variables: hello=" HI, WORLD ",key1="value1",key2="value2",key3="value3"
}

func TestMapVar(t *testing.T) {
	t.Run("should parse well formed flags", func(t *testing.T) {
		testcases := []struct {
			args     []string
			expected map[string]string
		}{
			{
				args: []string{`-m`, `HELLO=WORLD`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD,HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO="WORLD,HALLO=WELT"`},
				expected: map[string]string{
					`HELLO`: `WORLD,HALLO=WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO="WORLD",HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO=" HI, WORLD "`},
				expected: map[string]string{
					`HELLO`: ` HI, WORLD `,
				},
			}, {
				args: []string{`-m`, `HELLO = WORLD`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD`, `-m`, `HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD`, `-m`, `HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD`, `-m`, `HALLO=WELT,HELLO=world`},
				expected: map[string]string{
					`HELLO`: `world`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD,HELLO=world`},
				expected: map[string]string{
					`HELLO`: `world`,
				},
			}, {
				args: []string{`-m`, `HELLO=" WOR "	LD`},
				expected: map[string]string{
					`HELLO`: ` WOR LD`,
				},
			}, {}, {
				args: []string{`-m`, `HELLO=" WO "	"R"	LD " "`},
				expected: map[string]string{
					`HELLO`: ` WO RLD `,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD=HALLO`},
				expected: map[string]string{
					`HELLO`: `WORLD=HALLO`,
				},
			},
		}

		for _, tt := range testcases {
			mv := MapVar{}

			fs := flag.NewFlagSet("map", flag.ContinueOnError)
			fs.Var(mv, "m", "test")

			if err := fs.Parse(tt.args); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !maps.Equal(tt.expected, mv) {
				t.Errorf("unexpected parsed args: %v (%v)", mv, tt.args)
			}
		}
	})

	t.Run("should error for malformed flags", func(t *testing.T) {
		testcases := [][]string{
			{`-m`, `HELLO`},
			{`-m`, `HELLO,WORLD`},
			{`-m`, `HELLO=WORLD,`},
			{`-m`, `HELLO="WORLD`},
			{`-m`, `HELLO=WORLD"`},
			{`-m`, `"HELLO"=WORLD`},
			{`-m`, `HELLO=WORLD,hi,HALLO=WELT`},
			{`-m`, `,`},
		}

		for _, tt := range testcases {
			fs := flag.NewFlagSet("map", flag.ContinueOnError)
			fs.Var(MapVar{}, "m", "test")

			if err := fs.Parse(tt); err == nil {
				t.Errorf("expected error for malformed flags: %v", tt)
			}
		}
	})
}
