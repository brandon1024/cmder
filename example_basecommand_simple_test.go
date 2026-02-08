package cmder_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/brandon1024/cmder"
)

// This example demonstrates the simplest usage of [cmder.BaseCommand]. By using BaseCommand directly, you don't need to
// define your own command types. This can be a nice convenience for simple commands.
func ExampleBaseCommand() {
	exampleSetup()

	args := []string{"-Co-", "-"}

	if err := cmder.Execute(context.Background(), untar, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// 00000000  0a 27 75 6e 74 61 72 27  20 64 65 6d 6f 6e 73 74  |.'untar' demonst|
	// 00000010  72 61 74 65 73 20 74 68  65 20 73 69 6d 70 6c 65  |rates the simple|
	// 00000020  73 74 20 75 73 61 67 65  20 6f 66 20 5b 63 6d 64  |st usage of [cmd|
	// 00000030  65 72 2e 42 61 73 65 43  6f 6d 6d 61 6e 64 5d 2c  |er.BaseCommand],|
	// 00000040  20 75 73 69 6e 67 20 74  68 65 20 74 79 70 65 20  | using the type |
	// 00000050  74 6f 20 69 6d 70 6c 65  6d 65 6e 74 20 61 20 73  |to implement a s|
	// 00000060  69 6d 70 6c 65 20 63 6f  6d 6d 61 6e 64 20 74 68  |imple command th|
	// 00000070  61 74 20 72 65 61 64 73  0a 61 20 74 61 72 20 61  |at reads.a tar a|
	// 00000080  72 63 68 69 76 65 20 61  6e 64 20 64 75 6d 70 73  |rchive and dumps|
	// 00000090  20 69 6e 66 6c 61 74 65  64 20 63 6f 6e 74 65 6e  | inflated conten|
	// 000000a0  74 2e 0a                                          |t..|
}

const UntarDesc = `
'untar' demonstrates the simplest usage of [cmder.BaseCommand], using the type to implement a simple command that reads
a tar archive and dumps inflated content.
`

const UntarExamples = `
untar example.tar
untar -Co=example.out example.tar
untar -Co- - <example.tar
`

var (
	untar = &cmder.BaseCommand{
		CommandName: "untar",
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:     "untar [-o <file>] [-C] <file>",
			ShortHelp: "A simple demonstration of direct usage of BaseCommand.",
			Help:      UntarDesc,
			Examples:  UntarExamples,
		},
		RunFunc:       run,
		InitFlagsFunc: flags,
	}
)

var (
	output  string = "-"
	hexdump bool
)

var (
	in io.ReadWriter = os.Stdin
)

// Configure flags. We register two flags, '-o' and '-C'.
func flags(fs *flag.FlagSet) {
	fs.StringVar(&output, "o", output, "dump archive content to `file`")
	fs.BoolVar(&hexdump, "C", hexdump, "dump archive content in a canonical hex+ascii format")
}

// The command's run function.
func run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return cmder.ErrShowUsage
	}

	var out io.Writer = os.Stdout

	if args[0] != "-" {
		inputf, err := os.Open(args[0])
		if err != nil {
			return err
		}

		defer inputf.Close()
		in = inputf
	}

	if output != "-" {
		outputf, err := os.Create(output)
		if err != nil {
			return err
		}

		defer outputf.Close()
		out = outputf
	}

	if hexdump {
		dumper := hex.Dumper(out)
		defer dumper.Close()
		out = dumper
	}

	reader := tar.NewReader(in)
	for {
		_, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, reader); err != nil {
			return err
		}
	}

	return nil
}

// Setup for the example.
func exampleSetup() {
	in = &bytes.Buffer{}

	tw := tar.NewWriter(in)
	defer tw.Close()

	tw.WriteHeader(&tar.Header{
		Name: "usage-example",
		Mode: 0644,
		Size: int64(len(UntarDesc)),
	})
	tw.Write([]byte(UntarDesc))
}
