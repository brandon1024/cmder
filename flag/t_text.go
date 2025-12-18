package flag

import (
	"encoding"
	"fmt"
	"reflect"
)

var (
	// Compile-time check
	_ Getter = new(textT)
)

type textT struct {
	value encoding.TextUnmarshaler
}

func newTextT(v encoding.TextMarshaler, p encoding.TextUnmarshaler) *textT {
	ptrVal := reflect.ValueOf(p)
	if ptrVal.Kind() != reflect.Ptr {
		panic("variable p type must be a pointer")
	}

	defVal := reflect.ValueOf(v)
	if defVal.Kind() == reflect.Ptr {
		defVal = defVal.Elem()
	}

	if defVal.Type() != ptrVal.Type().Elem() {
		panic(fmt.Sprintf("default type does not match variable type: %v != %v", defVal.Type(), ptrVal.Type().Elem()))
	}

	ptrVal.Elem().Set(defVal)

	return &textT{p}
}

func (t textT) Set(v string) error {
	return t.value.UnmarshalText([]byte(v))
}

func (t textT) Get() any {
	return t.value
}

func (t textT) String() string {
	if m, ok := t.value.(encoding.TextMarshaler); ok {
		if b, err := m.MarshalText(); err == nil {
			return string(b)
		}
	}

	return ""
}
