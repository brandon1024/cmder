package getopt

import (
	"testing"
)

func TestIntVar(t *testing.T) {
	t.Run("should correctly parse integer values", func(t *testing.T) {
		var value int8

		i := Int(&value)

		if err := i.Set("127"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 127 {
			t.Fatalf("unexpected result: %d", value)
		}

		if err := i.Set("-128"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != -128 {
			t.Fatalf("unexpected result: %d", value)
		}

		if err := i.Set("0x0F"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 0x0F {
			t.Fatalf("unexpected result: %d", value)
		}
	})

	t.Run("should reject invalid integer values", func(t *testing.T) {
		var value int8

		i := Int(&value)

		if err := i.Set("128"); err == nil {
			t.Fatalf("expected error but was nil")
		}
		if err := i.Set("-129"); err == nil {
			t.Fatalf("expected error but was nil")
		}
	})

	t.Run("should correctly handle derived types", func(t *testing.T) {
		var value DerivedInt

		i := Int(&value)

		if err := i.Set("12345"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 12345 {
			t.Fatalf("unexpected result: %d", value)
		}
		if err := i.Set("-12345"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != -12345 {
			t.Fatalf("unexpected result: %d", value)
		}
	})
}

func TestUintVar(t *testing.T) {
	t.Run("should correctly parse integer values", func(t *testing.T) {
		var value uint8

		u := Uint(&value)

		if err := u.Set("255"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 255 {
			t.Fatalf("unexpected result: %d", value)
		}

		if err := u.Set("0"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 0 {
			t.Fatalf("unexpected result: %d", value)
		}

		if err := u.Set("0xFA"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 0xFA {
			t.Fatalf("unexpected result: %d", value)
		}
	})

	t.Run("should reject invalid integer values", func(t *testing.T) {
		var value uint8

		u := Uint(&value)

		if err := u.Set("256"); err == nil {
			t.Fatalf("expected error but was nil")
		}
		if err := u.Set("-1"); err == nil {
			t.Fatalf("expected error but was nil")
		}
	})

	t.Run("should correctly handle derived types", func(t *testing.T) {
		var value DerivedUint

		u := Uint(&value)

		if err := u.Set("12345"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 12345 {
			t.Fatalf("unexpected result: %d", value)
		}
	})
}

type (
	DerivedInt  int16
	DerivedUint uint16
)
