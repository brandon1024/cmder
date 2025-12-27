package flag

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestFlagSet(t *testing.T) {
	t.Run("Var", func(t *testing.T) {
		fs := NewFlagSet("test", ContinueOnError)

		t.Run("should panic if flag name begins with hyphen", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("no panic")
				}
			}()

			fs.String("-test", "", "string var")
		})

		t.Run("should panic if flag name ends with hyphen", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("no panic")
				}
			}()

			fs.String("test-", "", "string var")
		})

		t.Run("should panic if flag name contains invalid characters", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("no panic")
				}
			}()

			fs.String("test=test", "", "string var")
		})

		t.Run("should panic if flag already registered", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("no panic")
				}
			}()

			fs.String("test", "", "string var")
			fs.String("test", "", "string var")
		})

		t.Run("should update variable with default value", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			if output != "-" {
				t.Fatalf("output var not updated with expected default value")
			}
			if count != 12 {
				t.Fatalf("count var not updated with expected default value")
			}
		})

		t.Run("should allow flag names with periods", func(t *testing.T) {
			var addr string

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&addr, "web.listen-address", ":8080", "bind `address` for the web server")

			if addr != ":8080" {
				t.Fatalf("addr var not updated with expected default value")
			}
		})
	})

	t.Run("Lookup", func(t *testing.T) {
		fs := NewFlagSet("test", ContinueOnError)

		t.Run("should return nil if flag doesn't exist", func(t *testing.T) {
			if result := fs.Lookup("unknown"); result != nil {
				t.Fatalf("unexpected result: expected nil but was %v", result)
			}
		})

		t.Run("should return nil if flag doesn't exist", func(t *testing.T) {
			fs.String("output", "-", "output file")

			if result := fs.Lookup("output"); result == nil {
				t.Fatalf("unexpected result: nil")
			}
		})
	})

	t.Run("Parse", func(t *testing.T) {
		t.Run("should parse long flags", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--output", "test.out", "--count", "0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 0x12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
		})

		t.Run("should parse long flags with inline values", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--output=test.out", "--count=0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 0x12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
		})

		t.Run("should parse boolean long flags", func(t *testing.T) {
			var (
				b1 bool
				b2 bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.BoolVar(&b1, "b1", false, "boolean flag b1")
			fs.BoolVar(&b2, "b2", true, "boolean flag b1")

			err := fs.Parse([]string{"--b1", "--b2"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b1 != true {
				t.Fatalf("b1 var not updated with expected value: %v", b1)
			}
			if b2 != true {
				t.Fatalf("b2 var not updated with expected value: %v", b2)
			}
		})

		t.Run("should parse boolean long flags with inline values", func(t *testing.T) {
			var (
				b1 bool
				b2 bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.BoolVar(&b1, "b1", false, "boolean flag b1")
			fs.BoolVar(&b2, "b2", true, "boolean flag b1")

			err := fs.Parse([]string{"--b1=true", "--b2=false"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b1 != true {
				t.Fatalf("b1 var not updated with expected value: %v", b1)
			}
			if b2 != false {
				t.Fatalf("b2 var not updated with expected value: %v", b2)
			}
		})

		t.Run("should parse short flags", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")

			err := fs.Parse([]string{"-o", "test.out", "-c", "0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 0x12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
		})

		t.Run("should parse combined short flags", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-OC", "-otest.out", "-c0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 0x12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
			if b1 != true {
				t.Fatalf("b1 var not updated with expected value: %v", b1)
			}
			if b2 != true {
				t.Fatalf("b2 var not updated with expected value: %v", b2)
			}
		})

		t.Run("should parse combined short flags with inline value", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-OCotest.out"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
			if b1 != true {
				t.Fatalf("b1 var not updated with expected value: %v", b1)
			}
			if b2 != true {
				t.Fatalf("b2 var not updated with expected value: %v", b2)
			}
		})

		t.Run("should stop processing arguments after --", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-OC", "-otest.out", "--", "-c0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
			if b1 != true {
				t.Fatalf("b1 var not updated with expected value: %v", b1)
			}
			if b2 != true {
				t.Fatalf("b2 var not updated with expected value: %v", b2)
			}
			if fs.NArg() != 1 {
				t.Fatalf("unexpected number of unparsed args: %d", fs.NArg())
			}
			if a := fs.Args()[0]; a != "-c0x12" {
				t.Fatalf("unexpected unparsed arg: %s", a)
			}
		})

		t.Run("should halt with an error when unrecognized flag found", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-OC", "-otest.out", "-U", "-c0x12"})
			if err == nil {
				t.Fatalf("expected error but was nil")
			}
			if err.Error() != "flag '-U' does not exist" {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("should halt without an error when unrecognized arg found", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-OC", "-otest.out", "U", "-c0x12"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if fs.NArg() != 2 {
				t.Fatalf("unexpected number of unparsed args: %d", fs.NArg())
			}
			if a := fs.Args()[0]; a != "U" {
				t.Fatalf("unexpected unparsed arg: %s", a)
			}
			if a := fs.Args()[1]; a != "-c0x12" {
				t.Fatalf("unexpected unparsed arg: %s", a)
			}
		})

		t.Run("should treat single hyphen as an arg", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--output", "test.out", "-", "--count", "0x80"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
			if count != 12 {
				t.Fatalf("count var not updated with expected value: %d", count)
			}
			if fs.NArg() != 3 {
				t.Fatalf("unexpected number of remaining arguments: %d", fs.NArg())
			}
			if !slices.Equal(fs.Args(), []string{"-", "--count", "0x80"}) {
				t.Fatalf("unexpected args: %v", fs.Args())
			}
		})

		t.Run("should return ErrHelp if short help flag given", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"-h", "--output", "test.out", "--count", "0x80"})
			if !errors.Is(err, ErrHelp) {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "-" {
				t.Fatalf("output var parsed erroneously")
			}
			if count != 12 {
				t.Fatalf("count var parsed erroneously")
			}
		})

		t.Run("should return ErrHelp if long help flag given", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--help", "--output", "test.out", "--count", "0x80"})
			if !errors.Is(err, ErrHelp) {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "-" {
				t.Fatalf("output var parsed erroneously")
			}
			if count != 12 {
				t.Fatalf("count var parsed erroneously")
			}
		})

		t.Run("should not return ErrHelp if help flag given but user defined", func(t *testing.T) {
			var (
				output string
				count  uint
				help   bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")
			fs.BoolVar(&help, "help", false, "show help")

			err := fs.Parse([]string{"--help", "--output", "test.out", "--count", "0x80"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test.out" {
				t.Fatalf("output var parsed erroneously")
			}
			if count != 0x80 {
				t.Fatalf("count var parsed erroneously")
			}
			if !help {
				t.Fatalf("help var parsed erroneously")
			}
		})

		t.Run("should return error if unknown flag found", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--output", "test.out", "--all", "--count", "0x80"})
			if err == nil {
				t.Fatalf("expected error but was nil")
			}
			if !strings.Contains(err.Error(), "flag '--all' does not exist") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("should return error if unknown short flag found", func(t *testing.T) {
			var (
				output string
				count  uint
				b1     bool
				b2     bool
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")
			fs.BoolVar(&b1, "O", false, "assume output file")
			fs.BoolVar(&b2, "C", false, "assume count of results")

			err := fs.Parse([]string{"-otest.out", "-OaC", "-c0x12"})
			if err == nil {
				t.Fatalf("expected error but was nil")
			}
			if !strings.Contains(err.Error(), "flag '-a' does not exist") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("should return error if long flag arg missing", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output file")
			fs.UintVar(&count, "count", 12, "number of results")

			err := fs.Parse([]string{"--output", "test.out", "--count"})
			if err == nil {
				t.Fatalf("expected error but was nil")
			}
			if !strings.Contains(err.Error(), "missing argument to flag '--count'") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("should return error if short flag arg missing", func(t *testing.T) {
			var (
				output string
				count  uint
			)

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "o", "-", "output file")
			fs.UintVar(&count, "c", 12, "number of results")

			err := fs.Parse([]string{"-otest.out", "-c"})
			if err == nil {
				t.Fatalf("expected error but was nil")
			}
			if !strings.Contains(err.Error(), "missing argument to flag '-c'") {
				t.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("should support shorthand flag aliasing", func(t *testing.T) {
			var output string

			fs := NewFlagSet("test", ContinueOnError)
			fs.StringVar(&output, "output", "-", "output `file`")

			flg := fs.Lookup("output")
			fs.Var(flg.Value, "o", flg.Usage)

			err := fs.Parse([]string{"--output=test.out", "-o", "test-1.out"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != "test-1.out" {
				t.Fatalf("output var not updated with expected value: %s", output)
			}
		})
	})

	t.Run("PrintDefaults", func(t *testing.T) {
		t.Run("should render short flags correctly", func(t *testing.T) {
			var buf bytes.Buffer

			fs := NewFlagSet("test", ContinueOnError)
			fs.SetOutput(&buf)

			fs.Uint("c", 12, "number of results")

			fs.PrintDefaults()

			expected := "   -c <uint> (default 12)\n        number of results\n"
			if buf.String() != expected {
				t.Fatalf("unexpected usage string: '%s'", buf.String())
			}
		})

		t.Run("should render long flags correctly", func(t *testing.T) {
			var buf bytes.Buffer

			fs := NewFlagSet("test", ContinueOnError)
			fs.SetOutput(&buf)

			fs.Uint("count", 12, "number of results")

			fs.PrintDefaults()

			expected := "  --count <uint> (default 12)\n        number of results\n"
			if buf.String() != expected {
				t.Fatalf("unexpected usage string: '%s'", buf.String())
			}
		})

		t.Run("should pick out quoted argument name correctly", func(t *testing.T) {
			var buf bytes.Buffer

			fs := NewFlagSet("test", ContinueOnError)
			fs.SetOutput(&buf)

			fs.Uint("count", 12, "`number` of results")

			fs.PrintDefaults()

			expected := "  --count <number> (default 12)\n        number of results\n"
			if buf.String() != expected {
				t.Fatalf("unexpected usage string: '%s'", buf.String())
			}
		})

		t.Run("should render in lexical order", func(t *testing.T) {
			var buf bytes.Buffer

			fs := NewFlagSet("test", ContinueOnError)
			fs.SetOutput(&buf)

			fs.Uint("count", 12, "`number` of results")
			fs.Uint("c", 12, "`number` of results")
			fs.String("output", "-", "output `file`")
			fs.String("o", "-", "output `file`")
			fs.Bool("all", false, "show `all`")
			fs.Bool("a", false, "show `all`")

			fs.PrintDefaults()

			expected := `   -a (default false)
        show all
  --all (default false)
        show all
   -c <number> (default 12)
        number of results
  --count <number> (default 12)
        number of results
   -o <file> (default "-")
        output file
  --output <file> (default "-")
        output file
`
			if buf.String() != expected {
				t.Fatalf("unexpected usage string: '%s'", buf.String())
			}
		})
	})
}
