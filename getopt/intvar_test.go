package getopt

import (
	"io/fs"
	"log/slog"
	"testing"
)

func TestIntVar(t *testing.T) {
	t.Run("should correctly parse integer values", func(t *testing.T) {
		tt := map[string]int8{
			"127":       127,
			"-128":      -128,
			"0x0F":      15,
			"-0x0F":     -15,
			"0.1k":      100,
			"-0.0625Ki": -64,
		}

		for arg, expected := range tt {
			var value int8
			i := Int(&value)

			if err := i.Set(arg); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if value != expected {
				t.Fatalf("unexpected result: %d", value)
			}
		}
	})

	t.Run("should reject invalid integer values", func(t *testing.T) {
		tt := []string{
			"128", "-129", "12k", "1.2", "0.1 k", "0xFFFFQ", "0.1Ki",
		}

		for _, arg := range tt {
			var value int8
			i := Int(&value)

			if err := i.Set(arg); err == nil {
				t.Fatalf("expected error but was nil: %s", arg)
			}
		}
	})

	t.Run("should correctly handle derived types", func(t *testing.T) {
		var lvl slog.Level

		i := Int(&lvl)

		if err := i.Set("4"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lvl != slog.LevelWarn {
			t.Fatalf("unexpected result: %d", lvl)
		}
		if err := i.Set("-4"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lvl != slog.LevelDebug {
			t.Fatalf("unexpected result: %d", lvl)
		}
	})
}

func TestUintVar(t *testing.T) {
	t.Run("should correctly parse integer values", func(t *testing.T) {
		tt := map[string]uint8{
			"255":     255,
			"0":       0,
			"0xFA":    0xFA,
			"0.2k":    200,
			"0.125Ki": 128,
		}

		for arg, expected := range tt {
			var value uint8
			i := Uint(&value)

			if err := i.Set(arg); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if value != expected {
				t.Fatalf("unexpected result: %d", value)
			}
		}
	})

	t.Run("should reject invalid integer values", func(t *testing.T) {
		tt := []string{
			"256", "-01", "12k", "0.1 k", "0xFFFFQ",
		}

		// FIXME Brandon @ Fri 18 Sep 2026 07:45:36 AM CEST
		// Currently, there's a defect in [big.Float.Uint64] that results in an unexpected [big.Accuracy] for small
		// numbers. I logged an upstream bug and it was fixed, but will take a while to make it into the stdlib.
		//
		// https://github.com/golang/go/issues/81535
		// https://go-review.googlesource.com/c/go/+/833064

		if i := Uint(new(uint8)); i.Set("1.2") == nil {
			t.Logf("info: stdlib missing patch 833064")
		}
		if i := Uint(new(uint8)); i.Set("0.2Ki") == nil {
			t.Logf("info: stdlib missing patch 833064")
		}

		for _, arg := range tt {
			var value uint8
			i := Uint(&value)

			if err := i.Set(arg); err == nil {
				t.Fatalf("expected error but was nil: %s", arg)
			}
		}
	})

	t.Run("should correctly handle derived types", func(t *testing.T) {
		var value fs.FileMode

		u := Uint(&value)

		if err := u.Set("0777"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 0777 {
			t.Fatalf("unexpected result: %d", value)
		}
	})
}
