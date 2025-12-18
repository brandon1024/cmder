package flag

var (
	// Compile-time check
	_ Value = new(funcT)
)

type funcT func(string) error

func (f funcT) String() string {
	return ""
}

func (f funcT) Set(v string) error {
	return f(v)
}
