// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"testing"
)

var (
	testErrNotification = errors.New("test notification")
	testErrNotifier     = errors.New("test notifier error")
)

func TestDispatcher_Error(t *testing.T) {
	tests := []struct {
		name             string
		notifier         *testDispatcherNotifier
		notification     error
		wantErr          error
		wantNotification error
	}{
		{
			name:     "empty notification",
			notifier: &testDispatcherNotifier{},
		},
		{
			name:             "non-empty notification",
			notifier:         &testDispatcherNotifier{},
			notification:     testErrNotification,
			wantNotification: testErrNotification,
		},
		{
			name:     "emulate notifier error",
			notifier: &testDispatcherNotifier{emulateError: testErrNotifier},
			wantErr:  testErrNotifier,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dispatch := NewDispatchNotifier(tc.notifier)

			err := dispatch.Error(tc.notification)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
			if !errors.Is(tc.notifier.got, tc.wantNotification) {
				t.Errorf("unexpected notification '%v'; want '%v'", tc.notifier.got, tc.wantNotification)
			}
		})
	}
}

type testDispatcherNotifier struct {
	emulateError error
	got          error
}

func (tdn *testDispatcherNotifier) Error(err error) error {
	tdn.got = err
	return tdn.emulateError
}
