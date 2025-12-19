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
  --web.listen-address <string> (default :9090)
        address on which to expose metrics
  --web.telemetry-path <string> (default /metrics)
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
	cmd.fs.String("s", "", "serial number of the device (e.g. 10293894a)")
	cmd.fs.String("addr", "", "address and port of the device (e.g. 192.168.1.1:4567)")
	cmd.fs.String("a", "", "address and port of the device (e.g. 192.168.1.1:4567)")
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
