package getopt

import (
	"encoding/csv"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// MapVar is a [flag.Value] for flags that accept map values. MapVar also implements [flag.Getter].
//
// MapVar parses flag values which are key=value pairs. Multiple key=value pairs may be comma separated (e.g.
// key1=value1,key2=value2). Keys should be alphanumeric. Values can contain commas and whitespace if enclosed in double
// quotes.
//
//	key1=value1
//	key1=value1,key2=value2
//	"key1=value, 1","key2=value, 2"
//	key1=v=1,key2=v=2
type MapVar map[string]string

// String returns the map, formatted as a set of key-value pairs.
func (m MapVar) String() string {
	var entries []string

	for _, k := range slices.Sorted(maps.Keys(m)) {
		entries = append(entries, k+"="+m[k])
	}

	var builder strings.Builder

	w := csv.NewWriter(&builder)
	if err := w.Write(entries); err != nil {
		panic(err)
	}

	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}

	return strings.TrimSuffix(builder.String(), "\n")
}

// Set fulfills the [flag.Value] interface. The given value must be a set of key-value pairs.
func (m MapVar) Set(value string) error {
	r := csv.NewReader(strings.NewReader(value))

	pairs, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("getopt: malformed map value: %s", value)
	}
	if len(pairs) != 1 {
		return fmt.Errorf("getopt: malformed map value: %s", value)
	}

	for _, pair := range pairs[0] {
		k, v, _ := strings.Cut(pair, "=")
		m[k] = v
	}

	return nil
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a
// map[string]string.
func (m MapVar) Get() any {
	return map[string]string(m)
}
