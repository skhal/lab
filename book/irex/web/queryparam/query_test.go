// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queryparam_test

import (
	"net/http"
	"testing"

	"github.com/skhal/lab/book/irex/web/queryparam"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		want   string
		wantOk bool
	}{
		{
			name: "no query param",
			req:  mustNewRequest(t, "GET", "http://example.com/"),
		},
		{
			name:   "empty query param",
			req:    mustNewRequest(t, "GET", "http://example.com/?q="),
			wantOk: true,
		},
		{
			name:   "with query param",
			req:    mustNewRequest(t, "GET", "http://example.com/?q=test+value"),
			want:   "test value",
			wantOk: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := queryparam.Query(tc.req)

			if got != tc.want {
				t.Errorf("Query(%s) = %q, _; want %q", tc.req.URL, got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("Query(%s) = _, %v; want %v", tc.req.URL, ok, tc.wantOk)
			}
		})
	}
}

func mustNewRequest(t *testing.T, method, url string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
