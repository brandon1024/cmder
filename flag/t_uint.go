package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(uintT)
)

type uintT uint

func newUintT(v uint, p *uint) *uintT {
	*p = v
	return (*uintT)(p)
}

func (u *uintT) String() string {
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *uintT) Set(v string) error {
	updated, err := strconv.ParseUint(v, 0, 0)
	if err != nil {
		return err
	}

	*u = uintT(updated)
	return nil
}

func (u *uintT) Get() any {
	return uint(*u)
}
