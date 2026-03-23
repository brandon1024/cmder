package getopt

import (
	"flag"
	"testing"
)

func TestCounterVar(t *testing.T) {
	t.Run("should increment counter value correctly", func(t *testing.T) {
		var counter uint

		fs := NewPosixFlagSet("test", flag.ContinueOnError)
		fs.Var(Counter(&counter), "verbose", "increase verbosity")
		Alias(fs.FlagSet, "verbose", "v")

		err := fs.Parse([]string{"--verbose", "--verbose", "-v", "-vvv", "--verbose=false"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if counter != 6 {
			t.Fatalf("unexpected counter value: %d", counter)
		}
	})

	t.Run("should panic if nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("no panic")
			}
		}()

		var c *uint
		if err := Counter(c).Set("true"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
