package getopt

import "testing"

func TestStringsVar(t *testing.T) {
	t.Run("should not panic if calling String on nil value", func(t *testing.T) {
		var z StringsVar

		if result := z.String(); result != "" {
			t.Fatalf("unexpected result: %s", result)
		}
	})
}
