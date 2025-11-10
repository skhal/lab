// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/gotest/internal/test"
)

func TestNewEventID(t *testing.T) {
	tests := []struct {
		name  string
		event *test.Event
		want  test.EventID
	}{
		{
			name:  "no package",
			event: &test.Event{Test: "test"},
			want:  test.EventID("test"),
		},
		{
			name:  "with package",
			event: &test.Event{Package: "package", Test: "test"},
			want:  test.EventID("package/test"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := test.NewEventID(tc.event)

			if tc.want != got {
				t.Errorf("test.NewEventID(%s) = %s; want %s", (*printableEvent)(tc.event), got, tc.want)
			}
		})
	}
}

type printableEvent test.Event

func (e *printableEvent) String() string {
	return fmt.Sprintf("{package:%q test:%q}", e.Package, e.Test)
}
