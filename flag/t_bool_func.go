package flag

var (
	// Compile-time check
	_ Value = new(boolFuncT)
)

type boolFuncT func(string) error

func (f boolFuncT) String() string {
	return ""
}

func (f boolFuncT) Set(v string) error {
	return f(v)
}

func (f boolFuncT) IsBoolFlag() bool {
	return true
}
