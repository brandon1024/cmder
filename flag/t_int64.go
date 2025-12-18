package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(int64T)
)

type int64T int64

func newInt64T(v int64, p *int64) *int64T {
	*p = v
	return (*int64T)(p)
}

func (i *int64T) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func (i *int64T) Set(v string) error {
	updated, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		return err
	}

	*i = int64T(updated)
	return nil
}

func (i *int64T) Get() any {
	return int64(*i)
}
