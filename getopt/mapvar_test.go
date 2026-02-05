package getopt

import (
	"flag"
	"maps"
	"testing"
)

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
				args: []string{`-m`, `HELLO,WORLD`},
				expected: map[string]string{
					`HELLO`: ``,
					`WORLD`: ``,
				},
			}, {
				args: []string{`-m`, `HELLO=WORLD,HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `"HELLO=WORLD,HALLO=WELT"`},
				expected: map[string]string{
					`HELLO`: `WORLD,HALLO=WELT`,
				},
			}, {
				args: []string{`-m`, `"HELLO=WORLD",HALLO=WELT`},
				expected: map[string]string{
					`HELLO`: `WORLD`,
					`HALLO`: `WELT`,
				},
			}, {
				args: []string{`-m`, `"HELLO= HI, WORLD "`},
				expected: map[string]string{
					`HELLO`: ` HI, WORLD `,
				},
			}, {
				args: []string{`-m`, `HELLO = WORLD`},
				expected: map[string]string{
					`HELLO `: ` WORLD`,
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
				args: []string{`-m`, `"HELLO= WO 	R	LD  "`},
				expected: map[string]string{
					`HELLO`: ` WO 	R	LD  `,
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

			// try parsing again from the output of [MapVar.String]
			mv2 := MapVar{}

			fs = flag.NewFlagSet("map", flag.ContinueOnError)
			fs.Var(mv2, "m", "test")

			if err := fs.Parse([]string{"-m", mv.String()}); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !maps.Equal(mv, mv2) {
				t.Errorf("unexpected parsed args: %v (%v)", mv, tt.args)
			}
		}
	})

	t.Run("should error for malformed flags", func(t *testing.T) {
		testcases := [][]string{
			{`-m`, `HELLO="WORLD`},
			{`-m`, `HELLO=WORLD"`},
			{`-m`, `"HELLO"=WORLD`},
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
