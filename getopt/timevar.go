package getopt

import (
	"time"
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
