// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var ErrInternalServer = errors.New("Server error") // generic error

const pageTimeout = 10 * time.Millisecond // time to generate a page

func serveFile(w http.ResponseWriter, req *http.Request) {
	req, cancel := withTimeout(req, pageTimeout)
	defer cancel()
	res := handle(req)
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

func handle(req *http.Request) *response {
	res := &response{
		data: make(chan []byte, 1),
		err:  make(chan error, 1),
		done: make(chan struct{}),
	}
	go func() {
		defer close(res.data)
		defer close(res.err)
		defer close(res.done)
		data, err := renderMarkdown(req.Context(), req.URL.Path)
		if err != nil {
			res.err <- err
			return
		}
		res.data <- data
	}()
	return res
}

const extMarkdown = ".md"

func renderMarkdown(ctx context.Context, path string) ([]byte, error) {
	if extMarkdown != filepath.Ext(path) {
		return nil, fmt.Errorf("%s does not have a markdown extension", path)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fname := filepath.Join(pwd, path)
	fi, err := os.Stat(fname)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s", path)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return ToHTML(ctx, data)
}
