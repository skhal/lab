// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mdview runs an HTTP server to serve markdown file at localhost:8080/.
package mdview

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// PageTimeout is the maximum time it may take to generate a page.
var PageTimeout = 10 * time.Millisecond

var ErrInternalServer = errors.New("Server error") // generic error

// Run starts an HTTP server to serve the file at localhost:8080/.
func Run(filename string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		req, cancel := withTimeout(req, PageTimeout)
		defer cancel()
		res := handle(req, filename)
		select {
		case <-res.Done():
		case <-req.Context().Done():
			err := req.Context().Err()
			if err == nil {
				err = ErrInternalServer
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := res.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", res.Data())
	})
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return s.ListenAndServe()
}

func withTimeout(req *http.Request, timeout time.Duration) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	req = req.WithContext(ctx)
	return req, cancel
}

type response struct {
	data chan []byte
	err  chan error
	done chan struct{}
}

func (res *response) Data() []byte {
	data, ok := <-res.data
	if !ok {
		return nil
	}
	return data
}

func (res *response) Err() error {
	err, ok := <-res.err
	if !ok {
		return nil
	}
	return err
}

func (res *response) Done() <-chan struct{} {
	return res.done
}

func handle(req *http.Request, filename string) *response {
	res := &response{
		data: make(chan []byte, 1),
		err:  make(chan error, 1),
		done: make(chan struct{}),
	}
	go func() {
		defer close(res.data)
		defer close(res.err)
		defer close(res.done)
		data, err := renderMarkdown(req.Context(), filename)
		if err != nil {
			res.err <- err
			return
		}
		res.data <- data
	}()
	return res
}

func renderMarkdown(ctx context.Context, filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s", filename)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return ToHTML(ctx, data)
}
