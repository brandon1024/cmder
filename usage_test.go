package cmder

import (
	"bytes"
	"flag"
	"log/slog"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/brandon1024/cmder/getopt"
)

const desc = `cmder - build powerful command-line applications in Go

'cmder' is a simple and flexible library for building command-line interfaces in Go. If you're coming from Cobra and
have used it for any length of time, you have surely had your fair share of difficulties with the library. 'cmder' will
feel quite a bit more comfortable and easy to use, and the wide range of examples throughout the project should help
you get started.

'cmder' takes a very opinionated approach to building command-line interfaces. The library will help you define,
structure and execute your commands, but that's about it. 'cmder' embraces simplicity because sometimes, less is better.

To define a new command, simply define a type that implements the 'Command' interface. If you want your command to have
additional behavior like flags or subcommands, simply implement the appropriate interfaces.
`

const examples = `
test --addr <addr> --serial-number <num>
test --log.level <level>
test --poll-interval <sec> --web.disable-exporter-metrics
`

const ExpectedDefaultHelp = `cmder - build powerful command-line applications in Go

'cmder' is a simple and flexible library for building command-line interfaces in Go. If you're coming from Cobra and
have used it for any length of time, you have surely had your fair share of difficulties with the library. 'cmder' will
feel quite a bit more comfortable and easy to use, and the wide range of examples throughout the project should help
you get started.

'cmder' takes a very opinionated approach to building command-line interfaces. The library will help you define,
structure and execute your commands, but that's about it. 'cmder' embraces simplicity because sometimes, less is better.

To define a new command, simply define a type that implements the 'Command' interface. If you want your command to have
additional behavior like flags or subcommands, simply implement the appropriate interfaces.

Usage:
  test [subcommands] [flags] [args]

Examples:
  test --addr <addr> --serial-number <num>
  test --log.level <level>
  test --poll-interval <sec> --web.disable-exporter-metrics

Available Commands:
  child-1        First child subcommand for parent
  child-2        Second child subcommand for parent

Flags:
  -a <address>, --addr=<address>
      address and port of the device (e.g. 192.168.1.1:4567)

  -t <key=value>, --arg=<key=value> (default k=v)
      render template with arguments (key=value)

  -r <value>, --hosts=<value> (default hello,world)
      specify remote hosts (e.g. tcp://127.0.0.1)

  --reconnect-interval=<duration> (default 1m0s)
      interval between connection attempts (e.g. 1m)

  -s <serial>, --serial-number=<serial>
      serial number of the device (e.g. 10293894a)

  --web.disable-exporter-metrics
      exclude metrics about the exporter itself (go_*)

  --web.listen-address=<string> (default :9090)
      address on which to expose metrics

  --web.telemetry-path=<string> (default /metrics)
      path under which to expose metrics

Use "test [command] --help" for more information about a command.
`

func TestHelp(t *testing.T) {
	child1 := &BaseCommand{
		CommandName: "child-1",
		CommandDocumentation: CommandDocumentation{
			Usage:     "child-1 [flags] [args]",
			ShortHelp: "First child subcommand for parent",
			Help:      desc,
			Examples:  examples,
		},
	}
	child2 := &BaseCommand{
		CommandName: "child-2",
		CommandDocumentation: CommandDocumentation{
			Usage:     "child-2 [flags] [args]",
			ShortHelp: "Second child subcommand for parent",
			Help:      desc,
			Examples:  examples,
		},
	}

	parent := &BaseCommand{
		CommandName: "test",
		CommandDocumentation: CommandDocumentation{
			Usage:     "test [subcommands] [flags] [args]",
			ShortHelp: "Usage text generation test",
			Help:      desc,
			Examples:  examples,
		},
		Children: []Command{child1, child2},
	}

	cmd := command{
		Command: parent,
		fs:      flag.NewFlagSet("cmd", flag.ContinueOnError),
	}

	cmd.fs.String("serial-number", "", "`serial` number of the device (e.g. 10293894a)")
	getopt.Alias(cmd.fs, "serial-number", "s")
	cmd.fs.String("addr", "", "`address` and port of the device (e.g. 192.168.1.1:4567)")
	getopt.Alias(cmd.fs, "addr", "a")

	cmd.fs.Var(getopt.MapVar{"k": "v"}, "arg", "render template with arguments (`key=value`)")
	getopt.Alias(cmd.fs, "arg", "t")

	cmd.fs.Var(&getopt.StringsVar{"hello", "world"}, "hosts", "specify remote hosts (e.g. tcp://127.0.0.1)")
	getopt.Alias(cmd.fs, "hosts", "r")

	cmd.fs.Duration("poll-interval", time.Duration(0), "attempt to poll the device status more frequently than advertised")
	getopt.Hide(cmd.fs, "poll-interval")

	cmd.fs.Duration("reconnect-interval", time.Minute, "interval between connection attempts (e.g. 1m)")
	cmd.fs.String("web.listen-address", ":9090", "address on which to expose metrics")
	cmd.fs.String("web.telemetry-path", "/metrics", "path under which to expose metrics")
	cmd.fs.Bool("web.disable-exporter-metrics", false, "exclude metrics about the exporter itself (go_*)")

	t.Run("DefaultHelpTemplate", func(t *testing.T) {
		t.Run("should render correctly", func(t *testing.T) {
			var buf bytes.Buffer

			err := help(cmd, &ExecuteOptions{
				helpTemplate: DefaultHelpTemplate,
				outputWriter: &buf,
			})
			assert(t, nilerr(err))

			t.Logf("result:\n%s", buf.String())

			if diff := cmp.Diff(ExpectedDefaultHelp, buf.String()); diff != "" {
				t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("should not render hidden subcommands", func(t *testing.T) {
			parent.Children = append(parent.Children, &BaseCommand{
				CommandName: "child-3",
				CommandDocumentation: CommandDocumentation{
					Usage:     "child-3 [flags] [args]",
					ShortHelp: "Third (hidden) child subcommand for parent",
					Help:      desc,
					Examples:  examples,
					IsHidden:  true,
				},
			})

			var buf bytes.Buffer

			err := help(cmd, &ExecuteOptions{
				helpTemplate: DefaultHelpTemplate,
				outputWriter: &buf,
			})
			assert(t, nilerr(err))

			t.Logf("result:\n%s", buf.String())

			if diff := cmp.Diff(ExpectedDefaultHelp, buf.String()); diff != "" {
				t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
			}
		})
	})
}

const ExpectedDefaultUsage = `Usage:
  example [flags] [args]

Examples:
  example --bool-flag --int-flag=1 arg

Flags:
  --bool-func-non-zero
      bool func with non-zero default value

  --func-zero=<value>, --bool-func-zero
      func with zero default value

  --bool-non-zero (default true)
      bool with non-zero default value

  --bool-zero
      bool with zero default value

  --counter-non-zero (default 12)
      counter with non-zero default value

  --counter-zero
      counter with zero default value

  --duration-non-zero=<duration> (default 1s)
      duration with non-zero default value

  --duration-zero=<duration>
      duration with zero default value

  --float64-non-zero=<float> (default 1)
      float64 with non-zero default value

  --float64-zero=<float>
      float64 with zero default value

  --func-non-zero=<value>
      func with non-zero default value

  --int-non-zero=<int> (default 12)
      int with non-zero default value

  --int-zero=<int>
      int with zero default value

  --int64-non-zero=<int> (default 13)
      int64 with non-zero default value

  --int64-zero=<int>
      int64 with zero default value

  --map-non-zero=<value> (default k=v)
      map flag with non-zero default value

  --map-zero=<value>
      map flag with zero default value

  --neg-bool-non-zero (default false)
      negated bool with non-zero default value

  --neg-bool-zero
      negated bool with zero default value

  --string-non-zero=<string> (default test)
      string with non-zero default value

  --string-zero=<string>
      string with zero default value

  --strings-non-zero=<value> (default item)
      string slice flag with non-zero default value

  --strings-zero=<value>
      string slice flag with zero default value

  --text-non-zero=<value> (default ERROR)
      textvar with non-zero default value

  --text-zero=<value> (default INFO)
      textvar with zero default value

  --time-non-zero=<value> (default 1970-01-04T00:00:00Z)
      time flag with non-zero default value

  --time-zero=<value>
      time flag with zero default value

  --uint-non-zero=<uint> (default 14)
      uint with non-zero default value

  --uint-zero=<uint>
      uint with zero default value

  --uint64-non-zero=<uint> (default 15)
      uint64 with non-zero default value

  --uint64-zero=<uint>
      uint64 with zero default value
`

func TestUsage(t *testing.T) {
	t.Run("should correctly render default flag values", func(t *testing.T) {
		cmd := command{
			Command: &BaseCommand{
				CommandName: "example",
				CommandDocumentation: CommandDocumentation{
					Usage:     "example [flags] [args]",
					ShortHelp: "Example Command",
					Help:      "This is a simple example.",
					Examples:  "example --bool-flag --int-flag=1 arg",
				},
			},
			fs: flag.NewFlagSet("cmd", flag.ContinueOnError),
		}

		// native var types
		var (
			boolZero        bool
			boolNonZero     = true
			boolFuncZero    func(string) error
			boolFuncNonZero = func(string) error { return nil }
			durationZero    time.Duration
			durationNonZero = time.Second
			floatZero       float64
			floatNonZero    = 1.0
			funcZero        func(string) error
			funcNonZero     = func(string) error { return nil }
			intZero         int
			intNonZero      = 12
			int64Zero       int64
			int64NonZero    int64 = 13
			uintZero        uint
			uintNonZero     uint = 14
			uint64Zero      uint64
			uint64NonZero   uint64 = 15
			stringZero      string
			stringNonZero   = "test"
			textZero        slog.Level
			textNonZero     = slog.LevelError
		)
		cmd.fs.BoolVar(&boolZero, "bool-zero", boolZero, "bool with zero default value")
		cmd.fs.BoolVar(&boolNonZero, "bool-non-zero", boolNonZero, "bool with non-zero default value")
		cmd.fs.BoolFunc("bool-func-zero", "bool func with zero default value", boolFuncZero)
		cmd.fs.BoolFunc("bool-func-non-zero", "bool func with non-zero default value", boolFuncNonZero)
		cmd.fs.DurationVar(&durationZero, "duration-zero", durationZero, "duration with zero default value")
		cmd.fs.DurationVar(&durationNonZero, "duration-non-zero", durationNonZero, "duration with non-zero default value")
		cmd.fs.Float64Var(&floatZero, "float64-zero", floatZero, "float64 with zero default value")
		cmd.fs.Float64Var(&floatNonZero, "float64-non-zero", floatNonZero, "float64 with non-zero default value")
		cmd.fs.Func("func-zero", "func with zero default value", funcZero)
		cmd.fs.Func("func-non-zero", "func with non-zero default value", funcNonZero)
		cmd.fs.IntVar(&intZero, "int-zero", intZero, "int with zero default value")
		cmd.fs.IntVar(&intNonZero, "int-non-zero", intNonZero, "int with non-zero default value")
		cmd.fs.Int64Var(&int64Zero, "int64-zero", int64Zero, "int64 with zero default value")
		cmd.fs.Int64Var(&int64NonZero, "int64-non-zero", int64NonZero, "int64 with non-zero default value")
		cmd.fs.UintVar(&uintZero, "uint-zero", uintZero, "uint with zero default value")
		cmd.fs.UintVar(&uintNonZero, "uint-non-zero", uintNonZero, "uint with non-zero default value")
		cmd.fs.Uint64Var(&uint64Zero, "uint64-zero", uint64Zero, "uint64 with zero default value")
		cmd.fs.Uint64Var(&uint64NonZero, "uint64-non-zero", uint64NonZero, "uint64 with non-zero default value")
		cmd.fs.StringVar(&stringZero, "string-zero", stringZero, "string with zero default value")
		cmd.fs.StringVar(&stringNonZero, "string-non-zero", stringNonZero, "string with non-zero default value")
		cmd.fs.TextVar(&textZero, "text-zero", textZero, "textvar with zero default value")
		cmd.fs.TextVar(&textNonZero, "text-non-zero", textNonZero, "textvar with non-zero default value")

		// custom var types
		var (
			negatedBoolZero    bool
			negatedBoolNonZero = true
			counterZero        uint16
			counterNonZero     uint16 = 12
			hiddenZero         time.Time
			hiddenNonZero      = time.Now()
			mapZero            = map[string]string{}
			mapNonZero         = map[string]string{"k": "v"}
			stringsZero        []string
			stringsNonZero     = []string{"item"}
			timeZero           time.Time
			timeNonZero        = time.UnixMilli(1000 * 60 * 60 * 24 * 3).UTC()
		)
		cmd.fs.Var(getopt.NegatedBool(&negatedBoolZero), "neg-bool-zero", "negated bool with zero default value")
		cmd.fs.Var(getopt.NegatedBool(&negatedBoolNonZero), "neg-bool-non-zero", "negated bool with non-zero default value")
		cmd.fs.Var(getopt.Counter(&counterZero), "counter-zero", "counter with zero default value")
		cmd.fs.Var(getopt.Counter(&counterNonZero), "counter-non-zero", "counter with non-zero default value")
		cmd.fs.Var(&getopt.HiddenVar{Value: getopt.Time(&hiddenZero)}, "hidden-zero", "hidden flag with zero default value")
		cmd.fs.Var(&getopt.HiddenVar{Value: getopt.Time(&hiddenNonZero)}, "hidden-non-zero", "hidden flag with non-zero default value")
		cmd.fs.Var(getopt.Map(mapZero), "map-zero", "map flag with zero default value")
		cmd.fs.Var(getopt.Map(mapNonZero), "map-non-zero", "map flag with non-zero default value")
		cmd.fs.Var(getopt.Strings(&stringsZero), "strings-zero", "string slice flag with zero default value")
		cmd.fs.Var(getopt.Strings(&stringsNonZero), "strings-non-zero", "string slice flag with non-zero default value")
		cmd.fs.Var(getopt.Time(&timeZero), "time-zero", "time flag with zero default value")
		cmd.fs.Var(getopt.Time(&timeNonZero), "time-non-zero", "time flag with non-zero default value")

		var buf bytes.Buffer

		err := usage(cmd, &ExecuteOptions{
			usageTemplate: DefaultUsageTemplate,
			outputWriter:  &buf,
		})
		assert(t, nilerr(err))

		t.Logf("result:\n%s", buf.String())

		if diff := cmp.Diff(ExpectedDefaultUsage, buf.String()); diff != "" {
			t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFlags(t *testing.T) {
	cmd := command{
		Command: &BaseCommand{
			CommandName: "test",
		},
	}

	t.Run("should group bool flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Bool("all", false, "bool flag")
		getopt.Alias(cmd.fs, "all", "a")
		getopt.Alias(cmd.fs, "a", "l")
		cmd.fs.Bool("show", false, "bool flag")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["all"]
		assert(t, eq(true, ok))
		assert(t, eq(3, len(group)))
		assert(t, eq("a", group[0].Name))
		assert(t, eq("l", group[1].Name))
		assert(t, eq("all", group[2].Name))

		group, ok = groups["show"]
		assert(t, eq(true, ok))
		assert(t, eq(1, len(group)))
		assert(t, eq("show", group[0].Name))
	})

	t.Run("should group string flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.String("from", "HEAD^", "string flag")
		getopt.Alias(cmd.fs, "from", "b")
		getopt.Alias(cmd.fs, "from", "B")
		cmd.fs.String("to", "HEAD", "string flag")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["from"]
		assert(t, eq(true, ok))
		assert(t, eq(3, len(group)))
		assert(t, eq("B", group[0].Name))
		assert(t, eq("b", group[1].Name))
		assert(t, eq("from", group[2].Name))

		group, ok = groups["to"]
		assert(t, eq(true, ok))
		assert(t, eq(1, len(group)))
		assert(t, eq("to", group[0].Name))
	})

	t.Run("should group duration flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Duration("since", time.Minute, "duration flag")
		getopt.Alias(cmd.fs, "since", "s")
		getopt.Alias(cmd.fs, "since", "f")
		cmd.fs.Duration("until", time.Duration(0), "duration flag")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["since"]
		assert(t, eq(true, ok))
		assert(t, eq(3, len(group)))
		assert(t, eq("f", group[0].Name))
		assert(t, eq("s", group[1].Name))
		assert(t, eq("since", group[2].Name))

		group, ok = groups["until"]
		assert(t, eq(true, ok))
		assert(t, eq(1, len(group)))
		assert(t, eq("until", group[0].Name))
	})

	t.Run("should group float flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Float64("epsilon", 0.00001, "float64 flag")
		getopt.Alias(cmd.fs, "epsilon", "e")
		getopt.Alias(cmd.fs, "e", "ep")
		cmd.fs.Float64("gamma", 0.01, "float64 flag")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["epsilon"]
		assert(t, eq(true, ok))
		assert(t, eq(3, len(group)))
		assert(t, eq("e", group[0].Name))
		assert(t, eq("ep", group[1].Name))
		assert(t, eq("epsilon", group[2].Name))

		group, ok = groups["gamma"]
		assert(t, eq(true, ok))
		assert(t, eq(1, len(group)))
		assert(t, eq("gamma", group[0].Name))
	})

	t.Run("should group int flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Int("page", 0, "int flag")
		getopt.Alias(cmd.fs, "page", "p")
		cmd.fs.Int("count", 100, "int flag")
		getopt.Alias(cmd.fs, "count", "c")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["page"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("p", group[0].Name))
		assert(t, eq("page", group[1].Name))

		group, ok = groups["count"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("c", group[0].Name))
		assert(t, eq("count", group[1].Name))
	})

	t.Run("should group int64 flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Int64("page", 0, "int64 flag")
		getopt.Alias(cmd.fs, "page", "a")
		cmd.fs.Int64("count", 100, "int64 flag")
		getopt.Alias(cmd.fs, "count", "b")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["page"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("a", group[0].Name))
		assert(t, eq("page", group[1].Name))

		group, ok = groups["count"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("b", group[0].Name))
		assert(t, eq("count", group[1].Name))
	})

	t.Run("should group uint flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Uint("page", 0, "uint flag")
		getopt.Alias(cmd.fs, "page", "x")
		cmd.fs.Uint("count", 100, "uint flag")
		getopt.Alias(cmd.fs, "count", "y")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["page"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("x", group[0].Name))
		assert(t, eq("page", group[1].Name))

		group, ok = groups["count"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("y", group[0].Name))
		assert(t, eq("count", group[1].Name))
	})

	t.Run("should group uint64 flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Uint64("page", 0, "uint64 flag")
		getopt.Alias(cmd.fs, "page", "px")
		cmd.fs.Uint64("count", 100, "uint64 flag")
		getopt.Alias(cmd.fs, "count", "cx")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["page"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("px", group[0].Name))
		assert(t, eq("page", group[1].Name))

		group, ok = groups["count"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("cx", group[0].Name))
		assert(t, eq("count", group[1].Name))
	})

	t.Run("should group mapvar flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Var(getopt.MapVar{}, "arg", "mapvar flag")
		getopt.Alias(cmd.fs, "arg", "a")
		cmd.fs.Var(getopt.MapVar{}, "template", "mapvar flag")
		getopt.Alias(cmd.fs, "template", "t")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["arg"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("a", group[0].Name))
		assert(t, eq("arg", group[1].Name))

		group, ok = groups["template"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("t", group[0].Name))
		assert(t, eq("template", group[1].Name))
	})

	t.Run("should group func flags", func(t *testing.T) {
		fn1 := func(v string) error {
			return nil
		}
		fn2 := func(v string) error {
			return nil
		}

		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.BoolFunc("verbose", "bool func flag", fn1)
		getopt.Alias(cmd.fs, "verbose", "v")
		cmd.fs.Func("optimize", "func flag", fn2)
		getopt.Alias(cmd.fs, "optimize", "O")

		groups := flags(cmd)
		assert(t, eq(2, len(groups)))

		group, ok := groups["verbose"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("v", group[0].Name))
		assert(t, eq("verbose", group[1].Name))

		group, ok = groups["optimize"]
		assert(t, eq(true, ok))
		assert(t, eq(2, len(group)))
		assert(t, eq("O", group[0].Name))
		assert(t, eq("optimize", group[1].Name))
	})
}
