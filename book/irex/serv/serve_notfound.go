// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"html/template"
	"net/http"
	"net/url"
)

var notFoundTemplate = template.Must(template.New("index.html").ParseFS(efs, "static/404/index.html"))

func serveNotFound(w http.ResponseWriter, req *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set(headerContentType, contentTypeTextHTML)
	d := struct {
		URL *url.URL
	}{
		req.URL,
	}
	return notFoundTemplate.Execute(w, d)
}
