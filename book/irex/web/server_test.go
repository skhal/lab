// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web_test

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/skhal/lab/book/irex/web"
	labtesting "github.com/skhal/lab/go/tests"
)

var update = flag.Bool("update", false, "update golden files")

func TestServer_startStop(t *testing.T) {
	s := new(web.Server)
	if err := s.ListenAndServe("127.0.0.1:0"); err != nil {
		t.Fatalf("unexpected server start error '%v'", err)
	}

	if err := s.Shutdown(t.Context()); err != nil {
		t.Errorf("unexpected server shutdown error '%v'", err)
	}

	if err := s.Err(); err != nil {
		t.Errorf("unexpected server error '%v'", err)
	}
}

func mustNewRequest(t *testing.T, method, url string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("unexpected request error '%v'", err)
	}
	return req
}

type address string

func (a address) Url(url string) string {
	return fmt.Sprintf("http://%s%s", a, url)
}

func TestServer_serve(t *testing.T) {
	s := new(web.Server)
	if err := s.ListenAndServe("127.0.0.1:0"); err != nil {
		t.Fatalf("unexpected server start error '%v'", err)
	}
	defer func() {
		if err := s.Shutdown(t.Context()); err != nil {
			t.Errorf("unexpected server shutdown error '%v'", err)
		}
		if err := s.Err(); err != nil {
			t.Errorf("unexpected server error '%v'", err)
		}
	}()
	addr := address(s.Addr())
	tests := []struct {
		name     string
		req      *http.Request
		wantCode status
		wantErr  bool
		golden   labtesting.GoldenFile
	}{
		{
			name:     "get",
			req:      mustNewRequest(t, http.MethodGet, addr.Url("/")),
			wantCode: http.StatusOK,
			golden:   labtesting.GoldenFile("testdata/no_query.html"),
		},
		{
			name:     "get invalid command",
			req:      mustNewRequest(t, http.MethodGet, addr.Url("/?q=test")),
			wantCode: http.StatusInternalServerError,
			golden:   labtesting.GoldenFile("testdata/query_test.html"),
		},
		{
			name:     "get ping",
			req:      mustNewRequest(t, http.MethodGet, addr.Url("/?q=ping")),
			wantCode: http.StatusOK,
			golden:   labtesting.GoldenFile("testdata/query_ping.html"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{}

			res, err := client.Do(tc.req)

			switch {
			case err != nil:
				if !tc.wantErr {
					t.Errorf("unexpected error '%v'", err)
				}
				// do not check the response
				return
			default:
				if tc.wantErr {
					t.Error("want error")
				}
			}
			defer res.Body.Close()
			if got := status(res.StatusCode); got != tc.wantCode {
				t.Errorf("unexpected code %d (%[1]s), want %d (%[2]s)", got, tc.wantCode)
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("unexpected error '%v'", err)
			}
			if *update {
				tc.golden.Write(t, string(body))
			}
			if d := tc.golden.Diff(t, string(body)); d != "" {
				t.Errorf("response mismatch (-want +got):\n%s", d)
			}
		})
	}
}

type status int

func (s status) String() string {
	return http.StatusText(int(s))
}
