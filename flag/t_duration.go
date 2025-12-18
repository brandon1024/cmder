package flag

import "time"

var (
	// Compile-time check
	_ Getter = new(durationT)
)

type durationT time.Duration

func newDurationT(v time.Duration, p *time.Duration) *durationT {
	*p = v
	return (*durationT)(p)
}

func (d *durationT) String() string {
	return time.Duration(*d).String()
}

func (d *durationT) Set(v string) error {
	duration, err := time.ParseDuration(v)
	if err != nil {
		return err
	}

	*d = durationT(duration)
	return nil
}

func (d *durationT) Get() any {
	return time.Duration(*d)
}
