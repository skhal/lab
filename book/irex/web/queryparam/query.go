// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package queryparam provides API to access query parameters from URL.
package queryparam

import "net/http"

const queryParamQ = "q"

// Query returns &q= query parameter from the request's URL. There is a second
// boolean return parameter to indicate whether the query was set in the URL to
// distinguish between present but empty and missing query parameter.
func Query(req *http.Request) (string, bool) {
	values := req.URL.Query()
	if !values.Has(queryParamQ) {
		return "", false
	}
	return values.Get(queryParamQ), true
}
