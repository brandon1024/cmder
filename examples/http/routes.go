package main

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
)

var (
	ErrUnauthorized = errors.New("access denied: bad credentials")
)

func (c *ServerCommand) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.RedirectHandler("/index.html", http.StatusMovedPermanently))
	mux.HandleFunc("GET /index.html", c.route(c.renderIndexPage))

	return mux
}

func (c *ServerCommand) route(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("incoming http request from client", "method", r.Method, "addr", r.RemoteAddr, "uri", r.URL.Path)

		// auth
		u, p, ok := r.BasicAuth()
		if !c.noAuth && (!ok || c.basicAuth != u+":"+p) {
			slog.Warn("client request denied: missing or invalid credentials", "method", r.Method, "addr", r.RemoteAddr,
				"uri", r.URL.Path)

			w.Header().Set("WWW-Authenticate", "Basic realm=\"cmder\", charset=\"UTF-8\"")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !c.noAuth {
			slog.Info("client authenticated", "method", r.Method, "addr", r.RemoteAddr, "uri", r.URL.Path, "user", u)
		}

		// configure max body size
		if c.maxBodySize >= 0 {
			r.Body = http.MaxBytesReader(w, r.Body, c.maxBodySize)
		}

		h.ServeHTTP(w, r)
	})
}

func (c *ServerCommand) renderIndexPage(w http.ResponseWriter, r *http.Request) {
	u, _, ok := r.BasicAuth()
	if !ok {
		u = "anonymous"
	}

	err := template.Must(
		template.New("index.html").Parse(`
			<!doctype html>
			<html lang="en-US">
				<head>
					<meta charset="utf-8" />
					<meta name="viewport" content="width=device-width" />
					<title>Hello World!</title>
				</head>
				<body>
					Hello, {{.}}!
				</body>
			</html>
		`),
	).Execute(w, u)

	if err != nil {
		slog.Error("bug: failed to execute template", "route", "/index.html")
	}
}
