package cmder

import (
	"context"
	"flag"
	"fmt"
	"slices"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("interspersed", func(t *testing.T) {
		var (
			l0f0, l0f1 uint
			l1f0, l1f1 string
			l2f0, l2f1 int
		)

		var result []string

		cmd := &BaseCommand{
			CommandName: "l0",
			InitFlagsFunc: func(fs *flag.FlagSet) {
				fs.UintVar(&l0f0, "l0f0", l0f0, "l0f0")
				fs.UintVar(&l0f1, "l0f1", l0f1, "l0f1")
			},
			Children: []Command{
				&BaseCommand{
					CommandName: "l1",
					InitFlagsFunc: func(fs *flag.FlagSet) {
						fs.StringVar(&l1f0, "l1f0", l1f0, "l1f0")
						fs.StringVar(&l1f1, "l1f1", l1f1, "l1f1")
					},
					Children: []Command{
						&BaseCommand{
							CommandName: "l2",
							InitFlagsFunc: func(fs *flag.FlagSet) {
								fs.IntVar(&l2f0, "l2f0", l2f0, "l2f0")
								fs.IntVar(&l2f1, "l2f1", l2f1, "l2f1")
							},
							RunFunc: func(ctx context.Context, args []string) error {
								result = args
								return nil
							},
						},
					},
				},
			},
		}

		t.Run("should parse interspersed args", func(t *testing.T) {
			l0f0, l0f1, l1f0, l1f1, l2f0, l2f1 = 0, 0, "", "", 0, 0
			result = nil

			err := Execute(t.Context(), cmd, WithInterspersedArgs(), WithArgs([]string{
				"--l0f0", "255", "--l0f1=27",
				"l1", "--l1f0", "254", "--l1f1=26",
				"l2", "--l2f0=253", "000", "--l2f1", "25", "111", "--", "--l2f0=255",
			}))

			assert(t, nilerr(err))
			assert(t, eq(255, l0f0))
			assert(t, eq(27, l0f1))
			assert(t, eq("254", l1f0))
			assert(t, eq("26", l1f1))
			assert(t, eq(253, l2f0))
			assert(t, eq(25, l2f1))
			assert(t, match([]string{"000", "111", "--l2f0=255"}, result))
		})

		t.Run("should not parse interspersed by default", func(t *testing.T) {
			l0f0, l0f1, l1f0, l1f1, l2f0, l2f1 = 0, 0, "", "", 0, 0
			result = nil

			err := Execute(t.Context(), cmd, WithArgs([]string{
				"--l0f0", "255", "--l0f1=27",
				"l1", "--l1f0", "254", "--l1f1=26",
				"l2", "--l2f0=253", "000", "--l2f1", "25", "111", "--", "--l2f0=255",
			}))

			assert(t, nilerr(err))
			assert(t, eq(255, l0f0))
			assert(t, eq(27, l0f1))
			assert(t, eq("254", l1f0))
			assert(t, eq("26", l1f1))
			assert(t, eq(253, l2f0))
			assert(t, eq(0, l2f1))
			assert(t, match([]string{"000", "--l2f1", "25", "111", "--", "--l2f0=255"}, result))
		})
	})
}

type result struct {
	res bool
	msg string
}

func assert(t *testing.T, res result) {
	if !res.res {
		t.Fatalf("expectation failed: %s", res.msg)
	}
}

func eq[T comparable](expected, actual T) result {
	return result{expected == actual, fmt.Sprintf("values not equal: expected %v but was %v", expected, actual)}
}

func nilerr(err error) result {
	return result{err == nil, fmt.Sprintf("unexpected error: %v", err)}
}

func match[S ~[]E, E comparable](expected, actual S) result {
	return result{slices.Equal(expected, actual), fmt.Sprintf("slices not equal: expected %v but was %v", expected, actual)}
}
