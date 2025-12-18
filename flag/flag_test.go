package flag

import (
	"testing"
	"time"
)

func TestUnquoteUsage(t *testing.T) {
	t.Run("should pick out name from usage string", func(t *testing.T) {
		var output string
		fs := NewFlagSet("echo", ContinueOnError)
		fs.StringVar(&output, "output", "-", "output `file` location")

		flg := fs.Lookup("output")
		if flg == nil {
			t.Fatalf("could not find flag 'output'")
		}

		name, usage := UnquoteUsage(flg)
		if name != "file" {
			t.Fatalf("unexpected name: %s", name)
		}
		if usage != "output file location" {
			t.Fatalf("unexpected usage: %s", usage)
		}
	})

	t.Run("should infer name from type", func(t *testing.T) {
		fs := NewFlagSet("echo", ContinueOnError)
		fs.String("output", "-", "output file location")
		fs.Duration("since", time.Second, "time since")
		fs.Float64("epsilon", 0.000001, "smidge")
		fs.Int("page", 0, "page count")
		fs.Int64("count", 100, "item count")
		fs.Uint("limit", 10, "a limit")
		fs.Uint64("max-bytes", 1<<10, "data limit")
		fs.Bool("all", false, "show all")
		fs.BoolFunc("verbose", "show all", func(string) error {
			return nil
		})
		fs.Func("level", "set level", func(string) error {
			return nil
		})

		if name, _ := UnquoteUsage(fs.Lookup("output")); name != "string" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("since")); name != "duration" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("epsilon")); name != "float" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("page")); name != "int" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("count")); name != "int" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("limit")); name != "uint" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("max-bytes")); name != "uint" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("all")); name != "" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("verbose")); name != "" {
			t.Fatalf("unexpected name: %v", name)
		}
		if name, _ := UnquoteUsage(fs.Lookup("level")); name != "value" {
			t.Fatalf("unexpected name: %v", name)
		}
	})
}
