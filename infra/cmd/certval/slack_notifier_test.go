// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

func TestSlackNotifier_Nofiy(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		err     error
		wantErr error
	}{
		{
			name: "must be POST request",
			err:  errors.New("test error"),
		},
		{
			name: "server fails",
			handler: func(w http.ResponseWriter, req *http.Request) {
				http.Error(w, "test server error", http.StatusInternalServerError)
			},
			err:     errors.New("test error"),
			wantErr: ErrNotify,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handle := func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPost {
					http.Error(w, "allowed methods: POST", http.StatusMethodNotAllowed)
					return
				}
				if tc.handler != nil {
					tc.handler(w, req)
				}
			}
			server := httptest.NewTLSServer(http.HandlerFunc(handle))
			defer server.Close()
			slack := pb.Slack_builder{
				Webhook: new(server.URL),
			}.Build()
			notifier := NewSlackNotifier(slack)

			err := notifier.Notify(tc.err)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
		})
	}
}
