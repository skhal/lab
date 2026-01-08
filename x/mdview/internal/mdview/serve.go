// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var ErrInternalServer = errors.New("Server error") // generic error

func listenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveFile)
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	fmt.Println("start server at", addr)
	return s.ListenAndServe()
}

func serveFile(w http.ResponseWriter, req *http.Request) {
	data, err := renderMarkdown(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

const extMarkdown = ".md"

func renderMarkdown(path string) ([]byte, error) {
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
	return ToHTML(data)
}
