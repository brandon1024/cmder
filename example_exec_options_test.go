package cmder_test

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"hash"

	"github.com/brandon1024/cmder"
)

func ExampleWithInterspersedArgs() {
	args := []string{"string-1", "-a", "md5", "string-2", "-c10", "string-3"}

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithInterspersedArgs(),
	}

	if err := cmder.Execute(context.Background(), hasher, ops...); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// 0559406fc9a7b5704464c303ebbba64c
}

func ExampleWithRelaxedFlagParsing() {
	// note that shorthand '--al' is permitted for flag '--algo'
	args := []string{"--al", "md5", "relaxed-parsing"}

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithRelaxedFlagParsing(),
	}

	if err := cmder.Execute(context.Background(), hasher, ops...); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// 21db31e27ddc3aef918b031bd978fa78
}

func ExampleWithNamedTemplate() {
	args := []string{"--help"}

	// note the usage of "date", "version" and "suite" templates in the title header
	groffTemplate := `
.TH {{ .Command.Name }} 1 "{{ template "date" }}" "{{ template "version" }}" "{{ template "suite" }}"

.SH NAME

{{ trim .Command.Name }} - {{ trim .Command.ShortHelpText }}

.SH SYNOPSIS

{{ trim .Command.UsageLine }}

.SH DESCRIPTION

{{ trim .Command.HelpText }}

.SH OPTIONS

-a, --algo=<algorithm>
	Select hashing algorithm (md5, sha1, sha256).

-c, --rounds=<count>
	Number of hashing rounds.

.SH EXAMPLES

{{ trim .Command.ExampleText }}
`

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithHelpTemplate(groffTemplate),
		cmder.WithNamedTemplate("suite", "hashtools(1)"),
		cmder.WithNamedTemplate("version", "0.1.2"),
		cmder.WithNamedTemplate("date", "2006-01-02"),
	}

	err := cmder.Execute(context.Background(), hasher, ops...)
	if !errors.Is(err, cmder.ErrShowHelp) {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// .TH hash 1 "2006-01-02" "0.1.2" "hashtools(1)"
	//
	// .SH NAME
	//
	// hash - Simple demonstration of interspersed arg parsing.
	//
	// .SH SYNOPSIS
	//
	// hash [<str>...] [<flags>...]
	//
	// .SH DESCRIPTION
	//
	// 'hash' demonstrates how cmder can be configured to parse args with interspersed args and flags. The command generates
	// and prints a hash of the concatenated command args.
	//
	// .SH OPTIONS
	//
	// -a, --algo=<algorithm>
	// 	Select hashing algorithm (md5, sha1, sha256).
	//
	// -c, --rounds=<count>
	// 	Number of hashing rounds.
	//
	// .SH EXAMPLES
	//
	// # with interspersed args
	// hash string-1 -a md5 string-2 -c 10 string-3
	//
	// # without interspersed args
	// hash -a md5 -c 10 string-1 string-2 string-3
}

const HashDesc = `
'hash' demonstrates how cmder can be configured to parse args with interspersed args and flags. The command generates
and prints a hash of the concatenated command args.
`

const HashExamples = `
# with interspersed args
hash string-1 -a md5 string-2 -c 10 string-3

# without interspersed args
hash -a md5 -c 10 string-1 string-2 string-3
`

var (
	hasher = &Hasher{
		BaseCommand: cmder.BaseCommand{
			CommandName: "hash",
			CommandDocumentation: cmder.CommandDocumentation{
				Usage:     "hash [<str>...] [<flags>...]",
				ShortHelp: "Simple demonstration of interspersed arg parsing.",
				Help:      HashDesc,
				Examples:  HashExamples,
			},
		},
		algo:   "sha256",
		rounds: 1,
	}
)

type Hasher struct {
	cmder.BaseCommand

	algo   string
	rounds uint
}

func (h *Hasher) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&h.algo, "algo", h.algo, "select hashing algorithm (md5, sha1, sha256)")
	fs.StringVar(&h.algo, "a", h.algo, "select hashing algorithm (md5, sha1, sha256)")
	fs.UintVar(&h.rounds, "rounds", h.rounds, "number of hashing rounds")
	fs.UintVar(&h.rounds, "c", h.rounds, "number of hashing rounds")
}

func (h *Hasher) Run(ctx context.Context, args []string) error {
	algos := map[string]hash.Hash{
		"md5":    md5.New(),
		"sha1":   sha1.New(),
		"sha256": sha256.New(),
	}

	alg, ok := algos[h.algo]
	if !ok {
		return fmt.Errorf("no such algorithm: %s", h.algo)
	}

	for range h.rounds {
		for _, s := range args {
			alg.Write([]byte(s))
		}
	}

	fmt.Printf("%x\n", alg.Sum(nil))

	return nil
}
