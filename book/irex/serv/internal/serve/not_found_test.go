// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skhal/lab/book/irex/serv/internal/serve"
)

func TestNotFound(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		wantCode status
	}{
		{
			name:     "get /test",
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			wantCode: http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		w := httptest.NewRecorder()

		err := serve.NotFound(w, tc.req)

		if err != nil {
			t.Errorf("unexpected error '%v'", err)
		}
		if got := status(w.Code); got != tc.wantCode {
			t.Errorf("unexpected code %s, want %s", got, tc.wantCode)
		}
	}
}

type status int

func (s status) String() string {
	return http.StatusText(int(s))
}
