// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"embed"
	"errors"
	"io"
	"net/http"

	"github.com/skhal/lab/book/irex/web/queryparam"
)

const (
	headerContentType   = "Content-Type"
	contentTypeTextHTML = "text/html;charset=utf-8"
)

var (
	//go:embed static
	efs embed.FS
)

// Root serves main page.
func Root(w http.ResponseWriter, req *http.Request) error {
	_, ok := queryparam.Query(req)
	if !ok {
		return serveWelcomePage(w)
	}
	return errors.New("page for query parameter is under constructions")
}

func serveWelcomePage(w http.ResponseWriter) error {
	f, err := efs.Open("static/root/index.html")
	if err != nil {
		return err
	}
	defer f.Close()
	w.Header().Set(headerContentType, contentTypeTextHTML)
	_, err = io.Copy(w, f)
	return err
}
