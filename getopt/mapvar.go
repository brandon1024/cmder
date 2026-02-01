package getopt

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// MapVar is a [flag.Value] for flags that accept map values. MapVar also implements [flag.Getter].
//
// MapVar parses flag values which are key=value pairs. Multiple key=value pairs may be comma separated (e.g.
// key1=value1,key2=vlue2). Keys must be alphanumeric. Values can contain periods and whitespace if enclosed in double
// quotes. Unquoted whitespace will be trimmed.
//
//	key1=value1
//	key1=value1,key2=value2
//	key1 = "value, 1"
//	key1=v=1,key2=v=2
type MapVar map[string]string

// String returns the map, formatted as a set of key-value pairs:
func (m MapVar) String() string {
	var entries []string

	for _, k := range slices.Sorted(maps.Keys(m)) {
		entries = append(entries, k+"="+strconv.Quote(m[k]))
	}

	return strings.Join(entries, ",")
}

// Set fulfills the [flag.Value] interface. The given value must be a set of key-value pairs.
func (m MapVar) Set(value string) error {
	var (
		entries         = make(map[string]string)
		quoted, inValue bool
		key, val        string
	)

	for pos, c := range value {
		switch {
		case c == '"' && !inValue:
			return fmt.Errorf("illegal mapvar value at position %d (quoted key): %s", pos, value)
		case !quoted && c == '"':
			quoted = true
		case quoted && c == '"':
			quoted = false
		case quoted && inValue:
			val += string(c)
		case quoted && !inValue:
			return fmt.Errorf("illegal mapvar value at position %d (quoted key): %s", pos, value)
		case unicode.IsSpace(c):
		case !inValue && c == ',':
			return fmt.Errorf("illegal mapvar value at position %d (malformed pair missing value): %s", pos, value)
		case inValue && c == '=':
			val += string(c)
		case c == '=':
			inValue = true
		case c == ',':
			entries[key] = val
			key, val, inValue = "", "", false
		case inValue:
			val += string(c)
		case !inValue:
			key += string(c)
		default:
			panic(fmt.Errorf("bug: unhandled case: pos=%d, c=%v, quoted=%v, inValue=%v, key='%s', val='%s'",
				pos, c, quoted, inValue, key, val))
		}
	}

	if quoted {
		return fmt.Errorf("illegal mapvar value (quote mismatch): %s", value)
	}
	if !inValue {
		return fmt.Errorf("illegal mapvar value (malformed pair missing value): %s", value)
	}

	entries[key] = val

	for k, v := range entries {
		m[k] = v
	}

	return nil
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a
// map[string]string.
func (m MapVar) Get() any {
	return map[string]string(m)
}
