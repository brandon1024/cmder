package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(intT)
)

type intT int

func newIntT(v int, p *int) *intT {
	*p = v
	return (*intT)(p)
}

func (i *intT) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func (i *intT) Set(v string) error {
	updated, err := strconv.ParseInt(v, 0, 0)
	if err != nil {
		return err
	}

	*i = intT(updated)
	return nil
}

func (i *intT) Get() any {
	return int(*i)
}
