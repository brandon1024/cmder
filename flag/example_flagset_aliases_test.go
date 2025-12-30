package flag_test

import (
	"fmt"
	"time"

	"github.com/brandon1024/cmder/flag"
)

// Two flags with the same [Value] can be used to create shorthand aliases for longer flags. For example, you may want
// to provide a shorthand `-u` for the longer `--until` flag. Simply use [FlagSet.Lookup] to fetch the flag and
// [FlagSet.Var] to create the shorthand.
func ExampleFlagSet() {
	var since, until time.Duration

	fs := flag.NewFlagSet("alises", flag.ContinueOnError)
	fs.DurationVar(&since, "since", -time.Minute, "show items since")
	fs.Var(alias(fs.Lookup("since"), "s"))
	fs.DurationVar(&until, "until", time.Duration(0), "show items until")
	fs.Var(alias(fs.Lookup("until"), "u"))

	args := []string{
		"--since=-12m", "-u", "1m",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("since: %s\n", since.String())
	fmt.Printf("until: %s\n", until.String())
	// Output:
	// since: -12m0s
	// until: 1m0s
}

func alias(flg *flag.Flag, name string) (flag.Value, string, string) {
	return flg.Value, name, flg.Usage
}
