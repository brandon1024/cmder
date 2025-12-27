package flag_test

import (
	"fmt"
	"time"

	"github.com/brandon1024/cmder/flag"
)

// TimeVar is a [flag.Value] for flags that accept timestamps in [time.RFC3339] format. TimeVar also implements
// [Flag.Getter].
type TimeVar time.Time

// String returns the [time.RFC3339] representation of the timestamp flag.
func (t *TimeVar) String() string {
	return time.Time(*t).Format(time.RFC3339)
}

// Set fulfills the [flag.Value] interface. The given value must be a correctly formatted [time.RFC3339] timestamp.
func (t *TimeVar) Set(value string) error {
	tm, err := time.Parse(time.RFC3339, value)
	if err == nil {
		*t = TimeVar(tm)
	}

	return err
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a
// [time.Time].
func (t *TimeVar) Get() any {
	return time.Time(*t)
}

// You can define custom types implementing [flag.Value] to handle different types of flags, like timestamps, IP
// addresses, string maps or slices.
func ExampleValue() {
	var since TimeVar

	fs := flag.NewFlagSet("custom", flag.ContinueOnError)
	fs.Var(&since, "since", "show items since")

	args := []string{
		"--since", "2025-01-01T00:00:00Z",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("since: %s\n", since.String())

	// Output:
	// since: 2025-01-01T00:00:00Z
}
