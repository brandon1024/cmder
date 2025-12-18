package flag

import "strconv"

var (
	// Compile-time check
	_ Getter = new(float64T)
)

type float64T float64

func newFloat64T(v float64, p *float64) *float64T {
	*p = v
	return (*float64T)(p)
}

func (f *float64T) String() string {
	return strconv.FormatFloat(float64(*f), 'g', -1, 64)
}

func (f *float64T) Set(v string) error {
	updated, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}

	*f = float64T(updated)
	return nil
}

func (f *float64T) Get() any {
	return float64(*f)
}
