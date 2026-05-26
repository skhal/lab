// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"embed"
	"io"
	"net/http"
)

const (
	headerContentType   = "Content-Type"
	contentTypeTextHTML = "text/html;charset=utf-8"
)

var (
	//go:embed static
	efs embed.FS
)

func serveRoot(w http.ResponseWriter, req *http.Request) error {
	f, err := efs.Open("static/root/index.html")
	if err != nil {
		return err
	}
	defer f.Close()
	w.Header().Set(headerContentType, contentTypeTextHTML)
	_, err = io.Copy(w, f)
	return err
}
