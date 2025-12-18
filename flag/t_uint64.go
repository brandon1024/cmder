package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(uint64T)
)

type uint64T uint64

func newUint64T(v uint64, p *uint64) *uint64T {
	*p = v
	return (*uint64T)(p)
}

func (u *uint64T) String() string {
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *uint64T) Set(v string) error {
	updated, err := strconv.ParseUint(v, 0, 64)
	if err != nil {
		return err
	}

	*u = uint64T(updated)
	return nil
}

func (u *uint64T) Get() any {
	return uint64(*u)
}
