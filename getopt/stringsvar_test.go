package getopt

import (
	"flag"
	"testing"
)

func TestStringsVar(t *testing.T) {
	t.Run("should not panic if calling String on nil value", func(t *testing.T) {
		var z StringsVar

		if result := z.String(); result != "" {
			t.Fatalf("unexpected result: %s", result)
		}
	})

	t.Run("should permit empty string args", func(t *testing.T) {
		var ss []string

		fs := NewPosixFlagSet("strings", flag.ContinueOnError)
		fs.Var(Strings(&ss), "arg", "string slice arg")

		err := fs.Parse([]string{"--arg=", "--arg", ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(ss) != 2 {
			t.Fatalf("unexpected slice length: %d", len(ss))
		}
		if ss[0] != "" {
			t.Fatalf("unexpected slice element: %s", ss[0])
		}
		if ss[1] != "" {
			t.Fatalf("unexpected slice element: %s", ss[0])
		}
	})

	t.Run("should append to slice", func(t *testing.T) {
		var ss []string

		fs := NewPosixFlagSet("strings", flag.ContinueOnError)
		fs.Var(Strings(&ss), "arg", "string slice arg")

		err := fs.Parse([]string{"--arg=hello world", "--arg", "hallo welt"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(ss) != 2 {
			t.Fatalf("unexpected slice length: %d", len(ss))
		}
		if ss[0] != "hello world" {
			t.Fatalf("unexpected slice element: %s", ss[0])
		}
		if ss[1] != "hallo welt" {
			t.Fatalf("unexpected slice element: %s", ss[0])
		}
	})
}
