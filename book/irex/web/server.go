// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/skhal/lab/book/irex/web/internal/serve"
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

	// Address is the host:port to serve HTTP on.
	Address string
}

// Err returns the error from [http.Server.ListenAndServe] if any.
func (s *Server) Err() error {
	return s.err
}

// Serve starts HTTP server listening on [Server.Address] in a goroutine. Use
// [Server.Err] to check for errors.
func (s *Server) Serve() error {
	if s.httpServer != nil {
		return ErrServerRuns
	}
	h := s.handler()
	s.serve(h)
	return nil
}

func (s *Server) handler() http.Handler {
	wrap := func(h handlerFunc) http.HandlerFunc {
		return withTimeout(defaultHandleTimeout, handleError(h))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", wrap(serve.Root))
	return mux
}

func (s *Server) serve(h http.Handler) {
	s.httpServer = &http.Server{
		Addr:         s.Address,
		Handler:      h,
		ReadTimeout:  defaultServerReadTimeout,
		WriteTimeout: defaultServerWriteTimeout,
	}
	go func() {
		err := s.httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			s.err = err
		}
	}()
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
