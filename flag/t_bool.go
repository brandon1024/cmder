package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(boolT)
)

type boolT bool

func newBoolT(v bool, p *bool) *boolT {
	*p = v
	return (*boolT)(p)
}

func (b *boolT) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *boolT) Set(v string) error {
	updated, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}

	*b = boolT(updated)
	return nil
}

func (b *boolT) Get() any {
	return bool(*b)
}

func (b *boolT) IsBoolFlag() bool {
	return true
}
