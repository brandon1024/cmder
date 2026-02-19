# cmder

`cmder` is an opinionated library for building powerful command-line
applications in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/brandon1024/cmder.svg)](https://pkg.go.dev/github.com/brandon1024/cmder)
[![Go Report Card](https://goreportcard.com/badge/github.com/brandon1024/cmder)](https://goreportcard.com/report/github.com/brandon1024/cmder)

## Overview

`cmder` is a simple and flexible library for building command-line interfaces in
Go. `cmder` takes a very opinionated approach to building command-line
interfaces. The library will help you define, structure and execute your
commands, but that's about it. `cmder` embraces simplicity because sometimes,
less is better. The wide range of examples throughout the project should help
you get started.

To define a new command, simply define a type that implements the `Command`
interface. If you want your command to have additional behaviour like flags or
subcommands, simply implement the appropriate interfaces.

Here are some highlights:

- Bring your own types. `cmder` doens't force you to use special `command`
  structs. As long as you implement our narrow interfaces, you're good to go!
- `cmder` is unobtrustive. Define your command and execute it. Simplicity above
  all else!
- `cmder` is totally stateless making it super easy to unit test your commands.
  This isn't the case in other libraries.
- We take great pride in our documentation. If you find anything unclear, please
  let us know so we can fix it.

## Usage

First, include `cmder` in your project:

```shell
$ go get github.com/brandon1024/cmder
```

The easiest way to build commands is `cmder.BaseCommand`. For simple commands,
this is the cleanest way to go. This might feel a little familiar if you're
coming from [Cobra](https://github.com/spf13/cobra).

```go
package main

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

const HelloWorldHelpText = `hello-world - broadcast hello to the world

'hello-world' demonstrates how to build commands with the BaseCommand type.
`

const HelloWorldExamples = `
# broadcast hello to the world
hello-world from cmder
`

func run(ctx context.Context, args []string) error {
	fmt.Println("Hello World!")
	return nil
}

func main() {
	cmd := &cmder.BaseCommand{
		CommandName: "hello-world",
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:       "hello-world [<args>...]",
			ShortHelp:   "Simple demonstration of cmder",
			Help:        HelloWorldHelpText,
			Examples:    HelloWorldExamples,
		},
		RunFunc:     run,
	}

	if err := cmder.Execute(context.Background(), cmd); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
}
```

For more complex commands, you can define your own command type. By embedding
`cmder.BaseCommand`, your command automatically implements all of the important
interfaces needed to document your command, define flags, register subcommands,
and so on. You can then override the default behaviour with your own.

```go
package main

import (
	"context"
	"fmt"
	"flag"

	"github.com/brandon1024/cmder"
)

const BaseCommandExampleHelpText = `base-command - a simple example with struct embedding

'base-command' demonstrates how to build commands and subcommands with BaseCommand.
`

const BaseCommandExampleExamples = `
# broadcast hello to the world
base-command from cmder

# broadcast another message
base-command --msg 'hi bob!'
`

type BaseCommandExample struct {
	cmder.CommandDocumentation

	msg string
}

func (c *BaseCommandExample) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.msg, "m", "hello world", "message to broadcast")
	fs.StringVar(&c.msg, "msg", "hello world", "message to broadcast")
}

func (c *BaseCommandExample) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: %s\n", c.Name(), c.msg)
	return nil
}

func (c *BaseCommandExample) Name() string {
	return "base-command"
}

func main() {
	cmd := &BaseCommandExample{
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:       "base-command [-m | --msg <message>] [<args>...]",
			ShortHelp:   "A simple example with struct embedding",
			Help:        BaseCommandExampleHelpText,
			Examples:    BaseCommandExampleExamples,
		},
	}

	if err := cmder.Execute(context.Background(), cmd); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
		os.Exit(1)
	}
}
```

If you need even more flexibility, you can instead implement the interfaces that
are relevant for your command:

- `Command`: All commands and subcommands must implement this interface.
- `FlagInitializer`: If your command has flags, implement this interface.
- `Initializer`: If your command needs some initialization, implement this
  interface.
- `Destroyer`: If your command needs some teardown, implement this
  interface.

For more information, read through our package documentation on
[pkg.go.dev](https://pkg.go.dev/github.com/brandon1024/cmder).

## Development

To build the project and run tests:

```shell
$ make
```

## License

All software components herein are subject to the terms and conditions
specified in the [MIT License](./LICENSE).
