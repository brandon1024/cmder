package getopt

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// StringsVar is a [flag.Value] for flags that accept one or more string values. StringsVar also implements
// [flag.Getter].
//
// StringsVar collects string arguments into a slice. Multiple string values may be comma separated (e.g.
// value1,value2). Values containing commas may be enclosed in double quotes.
//
//	value
//	value1,value2
//	"value, 1","value, 2"
type StringsVar []string

// String returns the slice, formatted as comma-separated values.
func (s StringsVar) String() string {
	var builder strings.Builder

	w := csv.NewWriter(&builder)
	if err := w.Write([]string(s)); err != nil {
		panic(err)
	}

	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}

	return builder.String()
}

// Set fulfills the [flag.Value] interface.
func (s *StringsVar) Set(value string) error {
	r := csv.NewReader(strings.NewReader(value))

	values, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("getopt: malformed string slice value: %s", value)
	}
	if len(values) != 1 {
		return fmt.Errorf("getopt: malformed string slice value: %s", value)
	}

	for _, val := range values[0] {
		*s = append(*s, val)
	}

	return nil
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a
// []string.
func (s StringsVar) Get() any {
	return []string(s)
}
