package cmder

import (
	"bytes"
	"flag"
	"testing"
	"time"

	"github.com/brandon1024/cmder/getopt"
	"github.com/google/go-cmp/cmp"
)

const desc = `cmder - build powerful command-line applications in Go

'cmder' is a simple and flexible library for building command-line interfaces in Go. If you're coming from Cobra and
have used it for any length of time, you have surely had your fair share of difficulties with the library. 'cmder' will
feel quite a bit more comfortable and easy to use, and the wide range of examples throughout the project should help
you get started.

'cmder' takes a very opinionated approach to building command-line interfaces. The library will help you define,
structure and execute your commands, but that's about it. 'cmder' embraces simplicity because sometimes, less is better.

To define a new command, simply define a type that implements the 'Command' interface. If you want your command to have
additional behaviour like flags or subcommands, simply implement the appropriate interfaces.
`

const examples = `
test --addr <addr> --serial-number <num>
test --log.level <level>
test --poll-interval <sec> --web.disable-exporter-metrics
`

const ExpectedCobraUsageTemplate = `cmder - build powerful command-line applications in Go

'cmder' is a simple and flexible library for building command-line interfaces in Go. If you're coming from Cobra and
have used it for any length of time, you have surely had your fair share of difficulties with the library. 'cmder' will
feel quite a bit more comfortable and easy to use, and the wide range of examples throughout the project should help
you get started.

'cmder' takes a very opinionated approach to building command-line interfaces. The library will help you define,
structure and execute your commands, but that's about it. 'cmder' embraces simplicity because sometimes, less is better.

To define a new command, simply define a type that implements the 'Command' interface. If you want your command to have
additional behaviour like flags or subcommands, simply implement the appropriate interfaces.

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

  --web.disable-exporter-metrics (default false)
      exclude metrics about the exporter itself (go_*)

  --web.listen-address=<string> (default :9090)
      address on which to expose metrics

  --web.telemetry-path=<string> (default /metrics)
      path under which to expose metrics

Use "test [command] --help" for more information about a command.
`

const ExpectedStdFlagUsageTemplate = `usage: test [subcommands] [flags] [args]
  -a address
    	address and port of the device (e.g. 192.168.1.1:4567)
  -addr address
    	address and port of the device (e.g. 192.168.1.1:4567)
  -arg key=value
    	render template with arguments (key=value) (default k=v)
  -hosts value
    	specify remote hosts (e.g. tcp://127.0.0.1) (default hello,world)
  -poll-interval value
    	attempt to poll the device status more frequently than advertised (default 0s)
  -r value
    	specify remote hosts (e.g. tcp://127.0.0.1) (default hello,world)
  -reconnect-interval duration
    	interval between connection attempts (e.g. 1m) (default 1m0s)
  -s serial
    	serial number of the device (e.g. 10293894a)
  -serial-number serial
    	serial number of the device (e.g. 10293894a)
  -t key=value
    	render template with arguments (key=value) (default k=v)
  -web.disable-exporter-metrics
    	exclude metrics about the exporter itself (go_*)
  -web.listen-address string
    	address on which to expose metrics (default ":9090")
  -web.telemetry-path string
    	path under which to expose metrics (default "/metrics")
`

func TestUsage(t *testing.T) {
	child1 := &BaseCommand{
		CommandName: "child-1",
		Usage:       "child-1 [flags] [args]",
		ShortHelp:   "First child subcommand for parent",
		Help:        desc,
		Examples:    examples,
	}
	child2 := &BaseCommand{
		CommandName: "child-2",
		Usage:       "child-2 [flags] [args]",
		ShortHelp:   "Second child subcommand for parent",
		Help:        desc,
		Examples:    examples,
	}

	parent := &BaseCommand{
		CommandName: "test",
		Usage:       "test [subcommands] [flags] [args]",
		ShortHelp:   "Usage text generation test",
		Help:        desc,
		Examples:    examples,
		Children:    []Command{child1, child2},
	}

	cmd := command{
		Command: parent,
		fs:      flag.NewFlagSet("cmd", flag.ContinueOnError),
	}

	cmd.fs.String("serial-number", "", "`serial` number of the device (e.g. 10293894a)")
	cmd.fs.Var(alias(cmd.fs.Lookup("serial-number"), "s"))
	cmd.fs.String("addr", "", "`address` and port of the device (e.g. 192.168.1.1:4567)")
	cmd.fs.Var(alias(cmd.fs.Lookup("addr"), "a"))

	cmd.fs.Var(getopt.MapVar{"k": "v"}, "arg", "render template with arguments (`key=value`)")
	cmd.fs.Var(alias(cmd.fs.Lookup("arg"), "t"))

	cmd.fs.Var(&getopt.StringsVar{"hello", "world"}, "hosts", "specify remote hosts (e.g. tcp://127.0.0.1)")
	cmd.fs.Var(alias(cmd.fs.Lookup("hosts"), "r"))

	cmd.fs.Duration("poll-interval", time.Duration(0), "attempt to poll the device status more frequently than advertised")
	getopt.Hide(cmd.fs.Lookup("poll-interval"))

	cmd.fs.Duration("reconnect-interval", time.Minute, "interval between connection attempts (e.g. 1m)")
	cmd.fs.String("web.listen-address", ":9090", "address on which to expose metrics")
	cmd.fs.String("web.telemetry-path", "/metrics", "path under which to expose metrics")
	cmd.fs.Bool("web.disable-exporter-metrics", false, "exclude metrics about the exporter itself (go_*)")

	t.Run("CobraUsageTemplate", func(t *testing.T) {
		t.Run("should render correctly", func(t *testing.T) {
			var buf bytes.Buffer
			UsageOutputWriter = &buf
			UsageTemplate = CobraUsageTemplate

			err := usage(cmd)
			assert(t, nilerr(err))

			if diff := cmp.Diff(ExpectedCobraUsageTemplate, buf.String()); diff != "" {
				t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("should not render hidden subcommands", func(t *testing.T) {
			parent.Children = append(parent.Children, &BaseCommand{
				CommandName: "child-3",
				Usage:       "child-3 [flags] [args]",
				ShortHelp:   "Third (hidden) child subcommand for parent",
				Help:        desc,
				Examples:    examples,
				IsHidden:    true,
			})

			var buf bytes.Buffer
			UsageOutputWriter = &buf
			UsageTemplate = CobraUsageTemplate

			err := usage(cmd)
			assert(t, nilerr(err))

			if diff := cmp.Diff(ExpectedCobraUsageTemplate, buf.String()); diff != "" {
				t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
			}
		})
	})

	t.Run("StdFlagUsageTemplate", func(t *testing.T) {
		t.Run("should render correctly", func(t *testing.T) {
			var buf bytes.Buffer
			UsageOutputWriter = &buf
			UsageTemplate = StdFlagUsageTemplate

			err := usage(cmd)
			assert(t, nilerr(err))

			if diff := cmp.Diff(ExpectedStdFlagUsageTemplate, buf.String()); diff != "" {
				t.Fatalf("usage text mismatch (-want +got):\n%s", diff)
			}
		})
	})
}

func TestFlags(t *testing.T) {
	cmd := command{
		Command: &BaseCommand{
			CommandName: "test",
			Usage:       "test [flags] [args]",
			ShortHelp:   "Usage text generation test",
			Help:        desc,
			Examples:    examples,
		},
	}

	t.Run("should group bool flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Bool("all", false, "bool flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("all"), "a"))
		cmd.fs.Var(alias(cmd.fs.Lookup("a"), "l"))
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
		cmd.fs.Var(alias(cmd.fs.Lookup("from"), "b"))
		cmd.fs.Var(alias(cmd.fs.Lookup("from"), "B"))
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
		cmd.fs.Var(alias(cmd.fs.Lookup("since"), "s"))
		cmd.fs.Var(alias(cmd.fs.Lookup("since"), "f"))
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
		cmd.fs.Var(alias(cmd.fs.Lookup("epsilon"), "e"))
		cmd.fs.Var(alias(cmd.fs.Lookup("e"), "ep"))
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
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "p"))
		cmd.fs.Int("count", 100, "int flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "c"))

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
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "a"))
		cmd.fs.Int64("count", 100, "int64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "b"))

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
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "x"))
		cmd.fs.Uint("count", 100, "uint flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "y"))

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
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "px"))
		cmd.fs.Uint64("count", 100, "uint64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "cx"))

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
		cmd.fs.Var(alias(cmd.fs.Lookup("arg"), "a"))
		cmd.fs.Var(getopt.MapVar{}, "template", "mapvar flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("template"), "t"))

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
		cmd.fs.BoolFunc("verbose", "boolfunc flag", fn1)
		cmd.fs.Var(alias(cmd.fs.Lookup("verbose"), "v"))
		cmd.fs.Func("optimize", "func flag", fn2)
		cmd.fs.Var(alias(cmd.fs.Lookup("optimize"), "O"))

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

func alias(flg *flag.Flag, name string) (flag.Value, string, string) {
	return flg.Value, name, flg.Usage
}
