// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"embed"
	"net/http"

	"github.com/skhal/lab/book/irex/intent"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query"
	"github.com/skhal/lab/book/irex/render"
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
		return render.Render(pb.Page_builder{}.Build(), w, req)
	}
	queryIntent, err := query.Understand(q)
	if err != nil {
		return err
	}
	page, err := intent.Fulfill(queryIntent)
	if err != nil {
		return err
	}
	return render.Render(page, w, req)
}
