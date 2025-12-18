package flag

var (
	// Compile-time check
	_ Getter = new(stringT)
)

type stringT string

func newStringT(v string, p *string) *stringT {
	*p = v
	return (*stringT)(p)
}

func (s *stringT) String() string {
	return string(*s)
}

func (s *stringT) Set(v string) error {
	*s = stringT(v)
	return nil
}

func (s *stringT) Get() any {
	return string(*s)
}
