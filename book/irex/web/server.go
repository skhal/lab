// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/skhal/lab/book/irex/web/internal/serve"
)

var (
	// keep-sorted start
	defaultHandleTimeout      = 100 * time.Millisecond
	defaultServerReadTimeout  = 1 * time.Second
	defaultServerWriteTimeout = 1 * time.Second
	// keep-sorted end
)

// Server is the finance server to serve the results page, plots, etc.
type Server struct {
	// Address is the host:port to serve HTTP on.
	Address string
}

// Run starts HTTP server listening on [Server.Address].
func (s *Server) Run() error {
	h := s.handler()
	return s.serve(h)
}

func (s *Server) handler() http.Handler {
	wrap := func(h handlerFunc) http.HandlerFunc {
		return withTimeout(defaultHandleTimeout, handleError(h))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", wrap(serve.Root))
	mux.HandleFunc("/", wrap(serve.NotFound))
	return mux
}

func (s *Server) serve(h http.Handler) error {
	hs := &http.Server{
		Addr:         s.Address,
		Handler:      h,
		ReadTimeout:  defaultServerReadTimeout,
		WriteTimeout: defaultServerWriteTimeout,
	}
	return hs.ListenAndServe()
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
