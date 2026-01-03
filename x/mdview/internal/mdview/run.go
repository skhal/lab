// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mdview runs an HTTP server to serve markdown file at localhost:8080/.
package mdview

import (
	"fmt"
	"net/http"
	"os"
)

// Run starts an HTTP server to serve the file at localhost:8080/.
func Run(filename string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		data, err := os.ReadFile(filename)
		if err != nil {
			msg := fmt.Sprintf("Failed to open %s", filename)
			http.Error(w, msg, http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}
		html := ToHTML(data)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(html))
	})
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return s.ListenAndServe()
}
