package getopt

import "testing"

func TestTimeVar(t *testing.T) {
	t.Run("should not panic if calling String on nil value", func(t *testing.T) {
		var z TimeVar

		if result := z.String(); result != "0001-01-01T00:00:00Z" {
			t.Fatalf("unexpected result: %s", result)
		}
	})
}
