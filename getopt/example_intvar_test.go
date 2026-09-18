package getopt_test

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates usage of [getopt.UintVar]. UintVar makes it easy to parse unsigned integers of varying
// widths and types, and with suffixes like `Gi` and `k`.
func ExampleUintVar() {
	flagset := flag.NewFlagSet("uintvar", flag.ContinueOnError)

	var (
		mode              fs.FileMode
		sizeLimit         uint64
		transferRateLimit uint16
	)

	flagset.Var(getopt.Uint(&mode), "mode", "create files with `mode` at the destination")
	flagset.Var(getopt.Uint(&sizeLimit), "max-size", "don't transfer any file larger than `SIZE`")
	flagset.Var(getopt.Uint(&transferRateLimit), "bwlimit", "specify the max transfer `rate`, specified in units per second")

	flagset.Parse([]string{
		"--mode", "0640",
		"--max-size", "12.5Gi",
		"--bwlimit", "4Ki",
	})

	fmt.Printf("mode: '%v'\n", mode)
	fmt.Printf("size limit: %d\n", sizeLimit)
	fmt.Printf("tx limit: %d\n", transferRateLimit)
	// Output:
	// mode: '-rw-r-----'
	// size limit: 13421772800
	// tx limit: 4096
}

// This example demonstrates usage of [getopt.IntVar]. UintVar makes it easy to parse unsigned integers of varying
// widths and types, and with suffixes like `Gi` and `k`.
func ExampleIntVar() {
	flagset := flag.NewFlagSet("intvar", flag.ContinueOnError)

	var (
		lvl    slog.Level
		offset int64
		exit   int8
	)

	flagset.Var(getopt.Int(&lvl), "log-level", "configure logging at the given `level`")
	flagset.Var(getopt.Int(&offset), "offset", "specify offset from the beginning of the file, or from the end if negative")
	flagset.Var(getopt.Int(&exit), "exit-status", "exit with this status on failure")

	flagset.Parse([]string{
		"--log-level", "-4",
		"--offset", "1.5M",
		"--exit-status", "-12",
	})

	fmt.Printf("level: '%v'\n", lvl)
	fmt.Printf("offset: %d\n", offset)
	fmt.Printf("exit status: %d\n", exit)
	// Output:
	// level: 'DEBUG'
	// offset: 1500000
	// exit status: -12
}
