package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brandon1024/cmder"
)

func main() {
	cmd := &HelloCommand{
		subcommands: []cmder.Command{&WorldCommand{}},
	}

	if err := cmder.Execute(context.Background(), cmd); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
		os.Exit(1)
	}
}
