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

	"github.com/skhal/lab/book/irex/query"
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
	q, ok := queryparam.Query(req)
	if !ok {
		return serveWelcomePage(w)
	}
	_, err := query.Understand(q)
	if err != nil {
		return err
	}
	return errors.New("page for query parameter is under construction")
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
