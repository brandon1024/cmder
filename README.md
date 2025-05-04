# cmder

`cmder` is an opinionated library for building powerful command-line
applications in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/brandon1024/cmder.svg)](https://pkg.go.dev/github.com/brandon1024/cmder)
[![Go Report Card](https://goreportcard.com/badge/github.com/brandon1024/cmder)](https://goreportcard.com/report/github.com/brandon1024/cmder)

## Usage

First, include `cmder` in your project:

```shell
$ go get github.com/brandon1024/cmder
```

```go
package main

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

const HelloWorldHelpText = `
'hello-world' demonstrates how to build commands and subcommands with BaseCommand.
`

const HelloWorldExamples = `
# broadcast hello to the world
hello-world from cmder
`

type HelloWorldCommand struct {
	cmder.BaseCommand
}

func (c *HelloWorldCommand) Run(ctx context.Context, args []string) error {
	fmt.Println("Hello World!")
	return nil
}

func main() {
	cmd := &BaseCommandExample{
		BaseCommand: cmder.BaseCommand{
			CommandName: "hello-world",
			Usage:       "hello-world [<args>...]",
			ShortHelp:   "Simple demonstration of cmder",
			Help:        HelloWorldHelpText,
			Examples:    HelloWorldExamples,
		},
	}

	if err := cmder.Execute(context.Background(), cmd); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
}
```

## Development

To build the project and run tests:

```shell
$ make
```

## License

All software components herein are subject to the terms and conditions
specified in the [MIT License](./LICENSE).
