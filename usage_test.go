package cmder

import (
	"bytes"
	"testing"
	"time"

	"github.com/brandon1024/cmder/flag"
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

'cmder' also offers a flag package which is a drop-in replacement for the standard library package of the same name for
parsing POSIX/GNU style flags.
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

'cmder' also offers a flag package which is a drop-in replacement for the standard library package of the same name for
parsing POSIX/GNU style flags.

Usage:
  test [flags] [args]

Examples:
  test --addr <addr> --serial-number <num>
  test --log.level <level>
  test --poll-interval <sec> --web.disable-exporter-metrics

Flags:
  -a <string>, --addr=<string>
      address and port of the device (e.g. 192.168.1.1:4567)
  --poll-interval=<duration> (default 0s)
      attempt to poll the device status more frequently than advertised
  --reconnect-interval=<duration> (default 1m0s)
      interval between connection attempts (e.g. 1m)
  -s <string>, --serial-number=<string>
      serial number of the device (e.g. 10293894a)
  --web.disable-exporter-metrics (default false)
      exclude metrics about the exporter itself (go_*)
  --web.listen-address=<string> (default :9090)
      address on which to expose metrics
  --web.telemetry-path=<string> (default /metrics)
      path under which to expose metrics
`

const ExpectedStdFlagUsageTemplate = `usage: test [flags] [args]
   -a <string>
        address and port of the device (e.g. 192.168.1.1:4567)
  --addr <string>
        address and port of the device (e.g. 192.168.1.1:4567)
  --poll-interval <duration> (default 0s)
        attempt to poll the device status more frequently than advertised
  --reconnect-interval <duration> (default 1m0s)
        interval between connection attempts (e.g. 1m)
   -s <string>
        serial number of the device (e.g. 10293894a)
  --serial-number <string>
        serial number of the device (e.g. 10293894a)
  --web.disable-exporter-metrics (default false)
        exclude metrics about the exporter itself (go_*)
  --web.listen-address <string> (default ":9090")
        address on which to expose metrics
  --web.telemetry-path <string> (default "/metrics")
        path under which to expose metrics
`

func TestUsage(t *testing.T) {
	cmd := command{
		Command: &BaseCommand{
			CommandName: "test",
			Usage:       "test [flags] [args]",
			ShortHelp:   "Usage text generation test",
			Help:        desc,
			Examples:    examples,
		},
		fs: flag.NewFlagSet("cmd", flag.ContinueOnError),
	}

	cmd.fs.String("serial-number", "", "serial number of the device (e.g. 10293894a)")
	cmd.fs.Var(alias(cmd.fs.Lookup("serial-number"), "s"))
	cmd.fs.String("addr", "", "address and port of the device (e.g. 192.168.1.1:4567)")
	cmd.fs.Var(alias(cmd.fs.Lookup("addr"), "a"))

	cmd.fs.Duration("poll-interval", time.Duration(0), "attempt to poll the device status more frequently than advertised")
	cmd.fs.Duration("reconnect-interval", time.Minute, "interval between connection attempts (e.g. 1m)")
	cmd.fs.String("web.listen-address", ":9090", "address on which to expose metrics")
	cmd.fs.String("web.telemetry-path", "/metrics", "path under which to expose metrics")
	cmd.fs.Bool("web.disable-exporter-metrics", false, "exclude metrics about the exporter itself (go_*)")

	t.Run("CobraUsageTemplate", func(t *testing.T) {
		t.Run("should render correctly", func(t *testing.T) {
			var buf bytes.Buffer
			UsageOutputWriter = &buf
			UsageTemplate = CobraUsageTemplate

			if err := usage(cmd); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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

			if err := usage(cmd); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["all"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 3 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "a" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "l" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[2].Name != "all" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["show"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 1 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
	})

	t.Run("should group string flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.String("from", "HEAD^", "string flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("from"), "b"))
		cmd.fs.Var(alias(cmd.fs.Lookup("from"), "B"))
		cmd.fs.String("to", "HEAD", "string flag")

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["from"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 3 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "B" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "b" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[2].Name != "from" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["to"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 1 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
	})

	t.Run("should group duration flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Duration("since", time.Minute, "duration flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("since"), "s"))
		cmd.fs.Var(alias(cmd.fs.Lookup("since"), "f"))
		cmd.fs.Duration("until", time.Duration(0), "duration flag")

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["since"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 3 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "f" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "s" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[2].Name != "since" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["until"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 1 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
	})

	t.Run("should group float flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Float64("epsilon", 0.00001, "float64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("epsilon"), "e"))
		cmd.fs.Var(alias(cmd.fs.Lookup("e"), "ep"))
		cmd.fs.Float64("gamma", 0.01, "float64 flag")

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["epsilon"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 3 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "e" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "ep" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[2].Name != "epsilon" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["gamma"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 1 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
	})

	t.Run("should group int flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Int("page", 0, "int flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "p"))
		cmd.fs.Int("count", 100, "int flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "c"))

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["page"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "p" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "page" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["count"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "c" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "count" {
			t.Fatalf("unexpected sort order in flag group")
		}
	})

	t.Run("should group int64 flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Int64("page", 0, "int64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "a"))
		cmd.fs.Int64("count", 100, "int64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "b"))

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["page"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "a" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "page" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["count"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "b" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "count" {
			t.Fatalf("unexpected sort order in flag group")
		}
	})

	t.Run("should group uint flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Uint("page", 0, "uint flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "x"))
		cmd.fs.Uint("count", 100, "uint flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "y"))

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["page"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "x" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "page" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["count"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "y" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "count" {
			t.Fatalf("unexpected sort order in flag group")
		}
	})

	t.Run("should group uint64 flags", func(t *testing.T) {
		cmd.fs = flag.NewFlagSet("cmd", flag.ContinueOnError)
		cmd.fs.Uint64("page", 0, "uint64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("page"), "px"))
		cmd.fs.Uint64("count", 100, "uint64 flag")
		cmd.fs.Var(alias(cmd.fs.Lookup("count"), "cx"))

		groups := flags(cmd)
		if len(groups) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}

		group, ok := groups["page"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "px" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "page" {
			t.Fatalf("unexpected sort order in flag group")
		}

		group, ok = groups["count"]
		if !ok {
			t.Fatalf("no group found")
		}
		if len(group) != 2 {
			t.Fatalf("unexpected number of flag groups: %v", group)
		}
		if group[0].Name != "cx" {
			t.Fatalf("unexpected sort order in flag group")
		}
		if group[1].Name != "count" {
			t.Fatalf("unexpected sort order in flag group")
		}
	})

	t.Run("should not group func flags which are not comparable", func(t *testing.T) {
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
		if len(groups) != 4 {
			t.Fatalf("unexpected number of flag groups: %v", groups)
		}
	})
}

func alias(flg *flag.Flag, name string) (flag.Value, string, string) {
	return flg.Value, name, flg.Usage
}
