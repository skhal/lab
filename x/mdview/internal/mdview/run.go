// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mdview runs an HTTP server to serve markdown file at localhost:8080/.
package mdview

import (
	"flag"
)

const defaultAddr = "localhost:8080" // http server address

// Config provides configuration for Markdown server.
type Config struct {
	httpAddr string
}

func (cfg *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.httpAddr, "http", defaultAddr, "http address")
}

// Run starts an HTTP server to serve the file at localhost:8080/.
func Run(cfg *Config) error {
	return listenAndServe(cfg.httpAddr)
}
