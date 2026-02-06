package cmder_test

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"

	"github.com/brandon1024/cmder"
)

const usage = `hash [<str>...] [<flags>...]`

const short = `Simple demonstration of interspersed arg parsing.`

const desc = `
'hash' desmonstrates how cmder can be configured to parse args with interspersed args and flags. The command accepts one
or more strings and generates a hash of the result.
`

const examples = `
# with interspersed args
hash string-1 -a md5 string-2 -c 10 string-3

# without interspersed args
hash -a md5 -c 10 string-1 string-2 string-3
`

type hasher struct {
	cmder.BaseCommand

	algo   string
	rounds uint
}

var (
	cmd cmder.Command = &hasher{
		BaseCommand: cmder.BaseCommand{
			CommandName: "hash",
			Usage:       usage,
			ShortHelp:   short,
			Help:        desc,
			Examples:    examples,
		},
		algo:   "sha256",
		rounds: 1,
	}
)

func (h *hasher) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&h.algo, "algo", h.algo, "select hashing algorithm (md5, sha1, sha256)")
	fs.StringVar(&h.algo, "a", h.algo, "select hashing algorithm (md5, sha1, sha256)")
	fs.UintVar(&h.rounds, "rounds", h.rounds, "number of hashing rounds")
	fs.UintVar(&h.rounds, "c", h.rounds, "number of hashing rounds")
}

func (h *hasher) Run(ctx context.Context, args []string) error {
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

func ExampleWithInterspersedArgs() {
	args := []string{"string-1", "-a", "md5", "string-2", "-c10", "string-3"}

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithInterspersedArgs(),
	}

	if err := cmder.Execute(context.Background(), cmd, ops...); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}

	// Output:
	// 0559406fc9a7b5704464c303ebbba64c
}
