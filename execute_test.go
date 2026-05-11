package cmder

import (
	"context"
	"flag"
	"net/http"
	"testing"
	"time"

	"github.com/brandon1024/cmder/getopt"
	"github.com/brandon1024/cmder/internal/tutil"
)

func TestExecute(t *testing.T) {
	t.Run("interspersed", func(t *testing.T) {
		var (
			l0f0, l0f1 uint
			l1f0, l1f1 string
			l2f0, l2f1 int
		)

		var result []string

		cmd := &BaseCommand{
			CommandName: "l0",
			InitFlagsFunc: func(fs *flag.FlagSet) {
				fs.UintVar(&l0f0, "l0f0", l0f0, "l0f0")
				fs.UintVar(&l0f1, "l0f1", l0f1, "l0f1")
			},
			Children: []Command{
				&BaseCommand{
					CommandName: "l1",
					InitFlagsFunc: func(fs *flag.FlagSet) {
						fs.StringVar(&l1f0, "l1f0", l1f0, "l1f0")
						fs.StringVar(&l1f1, "l1f1", l1f1, "l1f1")
					},
					Children: []Command{
						&BaseCommand{
							CommandName: "l2",
							InitFlagsFunc: func(fs *flag.FlagSet) {
								fs.IntVar(&l2f0, "l2f0", l2f0, "l2f0")
								fs.IntVar(&l2f1, "l2f1", l2f1, "l2f1")
							},
							RunFunc: func(ctx context.Context, args []string) error {
								result = args
								return nil
							},
						},
					},
				},
			},
		}

		t.Run("should parse interspersed args", func(t *testing.T) {
			l0f0, l0f1, l1f0, l1f1, l2f0, l2f1 = 0, 0, "", "", 0, 0
			result = nil

			err := Execute(t.Context(), cmd, WithInterspersedArgs(), WithArgs([]string{
				"--l0f0", "255", "--l0f1=27",
				"l1", "--l1f0", "254", "--l1f1=26",
				"l2", "--l2f0=253", "000", "--l2f1", "25", "111", "--", "--l2f0=255",
			}))

			tutil.Assert(t, tutil.NilErr(err))
			tutil.Assert(t, tutil.Eq(255, l0f0))
			tutil.Assert(t, tutil.Eq(27, l0f1))
			tutil.Assert(t, tutil.Eq("254", l1f0))
			tutil.Assert(t, tutil.Eq("26", l1f1))
			tutil.Assert(t, tutil.Eq(253, l2f0))
			tutil.Assert(t, tutil.Eq(25, l2f1))
			tutil.Assert(t, tutil.Match([]string{"000", "111", "--l2f0=255"}, result))
		})

		t.Run("should not parse interspersed by default", func(t *testing.T) {
			l0f0, l0f1, l1f0, l1f1, l2f0, l2f1 = 0, 0, "", "", 0, 0
			result = nil

			err := Execute(t.Context(), cmd, WithArgs([]string{
				"--l0f0", "255", "--l0f1=27",
				"l1", "--l1f0", "254", "--l1f1=26",
				"l2", "--l2f0=253", "000", "--l2f1", "25", "111", "--", "--l2f0=255",
			}))

			tutil.Assert(t, tutil.NilErr(err))
			tutil.Assert(t, tutil.Eq(255, l0f0))
			tutil.Assert(t, tutil.Eq(27, l0f1))
			tutil.Assert(t, tutil.Eq("254", l1f0))
			tutil.Assert(t, tutil.Eq("26", l1f1))
			tutil.Assert(t, tutil.Eq(253, l2f0))
			tutil.Assert(t, tutil.Eq(0, l2f1))
			tutil.Assert(t, tutil.Match([]string{"000", "--l2f1", "25", "111", "--", "--l2f0=255"}, result))
		})
	})

	t.Run("native flags", func(t *testing.T) {
		var (
			addr           string
			readTimeout    time.Duration
			writeTimeout   time.Duration
			maxHeaderBytes int
			maxBodySize    int64
			basicAuth      string
			noAuth         bool
		)

		var args []string

		cmd := &BaseCommand{
			CommandName: "native-flags",
			InitFlagsFunc: func(fs *flag.FlagSet) {
				fs.StringVar(&addr, "http.bind-addr", ":8080", "bind address for the web server")
				fs.DurationVar(&readTimeout, "http.read-timeout", time.Duration(0), "read timeout for requests")
				fs.DurationVar(&writeTimeout, "http.write-timeout", time.Duration(0), "write timeout for responses")
				fs.IntVar(&maxHeaderBytes, "http.max-header-size", http.DefaultMaxHeaderBytes, "max permitted size of the headers in a request")
				fs.Int64Var(&maxBodySize, "http.max-body-size", 1<<26, "max permitted size of the headers in a request")
				fs.StringVar(&basicAuth, "http.auth-basic", "", "basic auth credentials (in format user:pass)")
				fs.BoolVar(&noAuth, "http.no-auth", false, "disable basic auth")

				getopt.Alias(fs, "http.bind-addr", "a")
				getopt.Alias(fs, "http.read-timeout", "r")
				getopt.Alias(fs, "http.write-timeout", "w")
				getopt.Alias(fs, "http.max-header-size", "h")
				getopt.Alias(fs, "http.max-body-size", "b")
				getopt.Alias(fs, "http.auth-basic", "C")
				getopt.Alias(fs, "http.no-auth", "E")
			},
			RunFunc: func(ctx context.Context, a []string) error {
				args = a
				return nil
			},
		}

		t.Run("should correctly parse flags in standard flag libs format", func(t *testing.T) {
			err := Execute(t.Context(), cmd, WithNativeFlags(), WithArgs([]string{
				"-http.bind-addr", "0.0.0.0:8000",
				"--http.read-timeout", "10s",
				"-http.write-timeout=5s",
				"-h", "8096",
				"-b=65536",
				"-http.auth-basic", "U:P",
				"--",
				"-http.no-auth", "true",
			}))

			tutil.Assert(t, tutil.NilErr(err))
			tutil.Assert(t, tutil.Eq("0.0.0.0:8000", addr))
			tutil.Assert(t, tutil.Eq(10*time.Second, readTimeout))
			tutil.Assert(t, tutil.Eq(5*time.Second, writeTimeout))
			tutil.Assert(t, tutil.Eq(8096, maxHeaderBytes))
			tutil.Assert(t, tutil.Eq(65536, maxBodySize))
			tutil.Assert(t, tutil.Eq("U:P", basicAuth))
			tutil.Assert(t, tutil.Eq(false, noAuth))
			tutil.Assert(t, tutil.Match([]string{"-http.no-auth", "true"}, args))
		})
	})

	t.Run("help flags", func(t *testing.T) {
		t.Run("should not register help flags if defined by command", func(t *testing.T) {
			var showHelp bool

			cmd := &BaseCommand{
				CommandName: "help-cmd",
				InitFlagsFunc: func(fs *flag.FlagSet) {
					fs.BoolVar(&showHelp, "h", showHelp, "show help")
					fs.BoolVar(&showHelp, "help", showHelp, "show help")
				},
			}

			err := Execute(t.Context(), cmd, WithArgs([]string{"--help"}))
			tutil.Assert(t, tutil.NilErr(err))
			tutil.Assert(t, tutil.Eq(true, showHelp))
		})

		t.Run("should return ErrShowHelp if help flags not defined by command", func(t *testing.T) {
			cmd := &BaseCommand{
				CommandName: "help-cmd",
			}

			err := Execute(t.Context(), cmd, WithArgs([]string{"--help"}))
			tutil.Assert(t, tutil.IsErr(err, ErrShowHelp))
		})

		t.Run("should return ErrShowUsage if help flags not defined by command", func(t *testing.T) {
			cmd := &BaseCommand{
				CommandName: "help-cmd",
			}

			err := Execute(t.Context(), cmd, WithArgs([]string{"-h"}))
			tutil.Assert(t, tutil.IsErr(err, ErrShowUsage))
		})
	})
}
