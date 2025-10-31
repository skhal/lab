// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/registry"
)

func TestParseQuestionIDs(t *testing.T) {
	tests := []struct {
		name    string
		qq      []string
		wantIDs []registry.QuestionID
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name:    "valid",
			qq:      []string{"1"},
			wantIDs: []registry.QuestionID{1},
		},
		{
			name:    "invalid",
			qq:      []string{"one"},
			wantErr: info.ErrQuestionID,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ids, err := info.ParseQuestionIDs(tc.qq)

			if diff := cmp.Diff(tc.wantIDs, ids, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("info.ParseQuestionIDs(%v) mismatch (-want, +got):\n%s", tc.qq, diff)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("info.ParseQuestionIDs(%v) = _, %s; want %s", tc.qq, err, tc.wantErr)
			}
		})
	}
}

func TestQuestionIDError_Is(t *testing.T) {
	err := &info.QuestionIDError{ID: "123"}

	if want := info.ErrQuestionID; !errors.Is(err, want) {
		t.Errorf("errors.Is(%#v, %#v) mismatch; want match", err, want)
	}
}

func TestMultiQuestionIDError_Is(t *testing.T) {
	err := &info.MultiQuestionIDError{IDs: []string{"123"}}

	if want := info.ErrQuestionID; !errors.Is(err, want) {
		t.Errorf("errors.Is(%#v, %#v) mismatch; want match", err, want)
	}
}
