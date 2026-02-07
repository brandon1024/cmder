package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/brandon1024/cmder"
)

const ServerCommandUsageLine = `server [<options>...]`

const ServerCommandShortHelpText = `A simple webserver built with cmder`

const ServerCommandHelpText = `server - a simple webserver built with cmder

'server' is a simple webserver built with cmder. This example helps demonstrate real-world usage of cmder. The webserver
simply serves a basic rendered webpage.

Without any options, 'server' starts a webserver on port 8080 protected with generated basic auth credentials. You can
configure a different bind address with the '--http.bind-addr' option. You can configure basic auth credentials with the
'--http.auth-basic' option. To disable auth, provide the '--http.no-auth' flag.
`

const ServerCommandExamples = `
# start the webserver with default (generated) credentials on port 8080
$ server

# start the webserver without auth and bind to loopback interface
$ server --http.no-auth --http.bind-addr 127.0.0.1:8080

# start the webserver with credentials
$ server --http.auth-basic ${USERNAME}:${PASSWORD}
`

type ServerCommand struct {
	// Bind address for the web server. By default, a bind address of `:8080` is used.
	//
	// See [http.Server].
	addr string

	// Read timeout for requests. By default, no timeout is imposed.
	//
	// See [http.Server].
	readTimeout time.Duration

	// Write timeout for responses. By default, no timeout is imposed.
	//
	// See [http.Server].
	writeTimeout time.Duration

	// Limit the number of header bytes.
	//
	// See [http.Server].
	maxHeaderBytes int

	// Limit the size of the request body.
	//
	// See [http.MaxBytesReader].
	maxBodySize int64

	// If set, requests are authorized if basic auth credentials are provided in the requests 'Authorization' header.
	// If unset, credentials are generated at startup.
	//
	// The value of this field is a username and password with format `user:pass`.
	basicAuth string

	// If configured, basic auth is disabled.
	noAuth bool
}

func (c *ServerCommand) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.addr, "http.bind-addr", ":8080", "bind address for the web server")
	fs.DurationVar(&c.readTimeout, "http.read-timeout", time.Duration(0), "read timeout for requests")
	fs.DurationVar(&c.writeTimeout, "http.write-timeout", time.Duration(0), "write timeout for responses")
	fs.IntVar(&c.maxHeaderBytes, "http.max-header-size", http.DefaultMaxHeaderBytes, "max permitted size of the headers in a request")
	fs.Int64Var(&c.maxBodySize, "http.max-body-size", 1<<26, "max permitted size of the headers in a request")
	fs.StringVar(&c.basicAuth, "http.auth-basic", "", "basic auth credentials (in format user:pass)")
	fs.BoolVar(&c.noAuth, "http.no-auth", false, "disable basic auth")
}

func (c *ServerCommand) Initialize(ctx context.Context, args []string) error {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "error: too many arguments: %v\n", args)
		return cmder.ErrShowUsage
	}

	if !c.noAuth && c.basicAuth == "" {
		var (
			user = "admin"
			pass = uuid.New().String()
		)

		slog.Info("no credentials configured: using generated basic auth credentials", "user", user, "pass", pass)

		c.basicAuth = user + ":" + pass
	}

	return nil
}

func (c *ServerCommand) Run(ctx context.Context, args []string) error {
	s := &http.Server{
		Addr:           c.addr,
		Handler:        c.routes(),
		ReadTimeout:    c.readTimeout,
		WriteTimeout:   c.writeTimeout,
		MaxHeaderBytes: c.maxHeaderBytes,
	}

	slog.Info("starting web server", "addr", c.addr)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := s.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown server", "err", err)
		}
	}()

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (c *ServerCommand) Destroy(ctx context.Context, args []string) error {
	return nil
}

func (c *ServerCommand) Name() string {
	return "server"
}

func (c *ServerCommand) UsageLine() string {
	return ServerCommandUsageLine
}

func (c *ServerCommand) ShortHelpText() string {
	return ServerCommandShortHelpText
}

func (c *ServerCommand) HelpText() string {
	return ServerCommandHelpText
}

func (c *ServerCommand) ExampleText() string {
	return ServerCommandExamples
}

func main() {
	cmd := &ServerCommand{}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := cmder.Execute(ctx, cmd); err != nil {
		fmt.Printf("unexpected error occurred: %v\n", err)
		os.Exit(1)
	}
}
