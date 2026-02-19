package getopt

import (
	"flag"
	"testing"
)

func TestAlias(t *testing.T) {
	t.Run("should panic if target flag does not exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("no panic")
			}
		}()

		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		Alias(fs, "non-existent", "q")
	})

	t.Run("should register alias successfully", func(t *testing.T) {
		var quiet bool
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.BoolVar(&quiet, "quiet", quiet, "silence the cat")

		Alias(fs, "quiet", "q")

		if err := fs.Parse([]string{"-q", "true"}); err != nil {
			t.Fatalf("failed to parse flags: %v", err)
		}
		if !quiet {
			t.Fatalf("alias not triggered")
		}
	})
}
