// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"compress/gzip"
	"context"
	"embed"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/lpar/gzipped/v2"
	"github.com/skhal/lab/book/irex/intent"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query"
	"github.com/skhal/lab/book/irex/render"
	"github.com/skhal/lab/book/irex/web/queryparam"
)

// ErrServerRuns means the web server already runs (see [Server.Serve]).
var ErrServerRuns = errors.New("web server runs")

var (
	// keep-sorted start
	defaultHandleTimeout      = 100 * time.Millisecond
	defaultServerReadTimeout  = 1 * time.Second
	defaultServerWriteTimeout = 1 * time.Second
	// keep-sorted end
)

// Server is the finance server to serve the results page, plots, etc.
type Server struct {
	err        error
	httpServer *http.Server
	addr       string
}

// Addr returns the address the server listens on.
func (s *Server) Addr() string {
	return s.addr
}

// Err returns the error from [http.Server.ListenAndServe] if any.
func (s *Server) Err() error {
	return s.err
}

// ListenAndServe starts HTTP server listening on the addr in a goroutine. It
// pickss up a random port of the address does not have have or set to 0.
//
// It returns an error if it fails to bind to the address. Use [Server.Err] to
// check for server errors.
func (s *Server) ListenAndServe(addr string) error {
	if s.httpServer != nil {
		return ErrServerRuns
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.addr = l.Addr().String()
	s.serve(l)
	return nil
}

func (s *Server) serve(l net.Listener) {
	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      s.handler(),
		ReadTimeout:  defaultServerReadTimeout,
		WriteTimeout: defaultServerWriteTimeout,
	}
	go func() {
		err := s.httpServer.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			s.err = err
		}
	}()
}

//go:embed static
var staticFS embed.FS

func (s *Server) handler() http.Handler {
	wrap := func(h handlerFunc) http.HandlerFunc {
		return withTimeout(defaultHandleTimeout, compress(handleError(h)))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", wrap(s.handleRoot))
	mux.Handle("GET /static/", gzipped.FileServer(gzipped.FS(staticFS)))
	return mux
}

// Shutdown gracefully shuts down the server. It returns the error from
// [http.Server.Shutdown].
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func withTimeout(timeout time.Duration, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()
		h(w, req.WithContext(ctx))
	}
}

type handlerFunc func(w http.ResponseWriter, req *http.Request) error

// handleError reports errors from handlerFunc as internal server errors (x500).
// It also sets up a context with timeout.
func handleError(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
	}
}

func compress(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !acceptEncoding(req) {
			h(w, req)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		h(gzw, req)
	}
}

func acceptEncoding(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write implements [io.Writer] interface.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (s *Server) handleRoot(w http.ResponseWriter, req *http.Request) error {
	q, ok := queryparam.Query(req)
	if !ok {
		return render.Render(pb.Page_builder{}.Build(), w, req)
	}
	queryIntent, err := query.Understand(q)
	if err != nil {
		return err
	}
	page, err := intent.Fulfill(queryIntent)
	if err != nil {
		return err
	}
	return render.Render(page, w, req)
}
