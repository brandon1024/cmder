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
	"github.com/brandon1024/cmder/getopt"
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
	args := []string{"-h"}

	// override 'section.header' from the default usage template
	header := `
		{{- define "section.header" -}}
			{{- println "hash - hash text input" -}}
			{{- println -}}
		{{- end -}}
	`

	// custom template, later named 'section.footer'
	footer := `
		{{- println -}}
		{{- println "This tool is freely licensed under a permissive MIT license." -}}
	`

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithNamedTemplate("welcome", header),
		cmder.WithNamedTemplate("section.examples",
			`{{- printf "\nExamples:\n" }}{{- printf "  hash string-1 -a md5 string-2 -c 10 string-3\n" -}}`),
		cmder.WithNamedTemplate("section.footer", footer),
	}

	err := cmder.Execute(context.Background(), hasher, ops...)
	if !errors.Is(err, cmder.ErrShowUsage) {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// hash - hash text input
	//
	// Usage:
	//   hash [<str>...] [<flags>...]
	//
	// Examples:
	//   hash string-1 -a md5 string-2 -c 10 string-3
	//
	// Flags:
	//   -a <string>, --algo=<string> (default md5)
	//       select hashing algorithm (md5, sha1, sha256)
	//
	//   -h
	//       show command usage information
	//
	//   --help
	//       show command help information
	//
	//   -c <uint>, --rounds=<uint> (default 10)
	//       number of hashing rounds
	//
	// This tool is freely licensed under a permissive MIT license.
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
	fs.UintVar(&h.rounds, "rounds", h.rounds, "number of hashing rounds")

	getopt.Alias(fs, "algo", "a")
	getopt.Alias(fs, "rounds", "c")
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
