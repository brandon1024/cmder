package getopt_test

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.Counter], configuring the log level of the default logger.
func ExampleCounter() {
	var verbosity slog.Level

	fs := getopt.NewPosixFlagSet("counter", flag.ContinueOnError)
	fs.Var(getopt.Counter(&verbosity), "v", "increase verbosity")

	if err := fs.Parse([]string{"-vvv"}); err != nil {
		panic(err)
	}

	lvl := slog.LevelError - verbosity*4

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	})))

	fmt.Printf("verbosity: %s\n", lvl.String())
	// Output:
	// verbosity: DEBUG
}
