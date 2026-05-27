// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skhal/lab/book/irex/web/internal/serve"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		wantCode status
	}{
		{
			name:     "get",
			req:      httptest.NewRequest(http.MethodGet, "/", nil),
			wantCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		w := httptest.NewRecorder()

		err := serve.Root(w, tc.req)

		if err != nil {
			t.Errorf("unexpected error '%v'", err)
		}
		if got := status(w.Code); got != tc.wantCode {
			t.Errorf("unexpected code %s, want %s", got, tc.wantCode)
		}
	}
}
