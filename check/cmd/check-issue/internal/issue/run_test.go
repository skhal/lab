// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package issue_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/check/cmd/check-issue/internal/issue"
)

var TestErr = errors.New("test error")

func TestCheck(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: issue.ErrNoIssue,
		},
		{
			name:    "no tag",
			s:       "Test data",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "no issue tag",
			s:    "NO_ISSUE",
		},
		{
			// can't mix the issue with other text on the line
			name:    "no issue tag with prefix",
			s:       "prefx NO_ISSUE",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "no issue tag with suffix",
			s:       "NO_ISSUE suffix",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "no issue tag with description",
			s:    "NO_ISSUE: N/A",
		},
		{
			// can't mix the issue with other text on the line
			name:    "no issue tag with description and prefix",
			s:       "prefix NO_ISSUE: N/A",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "local issue",
			s:    "Issue #123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "local issue with prefix",
			s:       "prefix Issue #123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "local issue with suffix",
			s:       "Issue #123 suffix",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "close local issue",
			s:    "Close #123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "close local issue with prefix",
			s:       "prefix Close #123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "close local issue with suffix",
			s:       "Close #123 suffix",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "fix local issue",
			s:    "Fix #123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "fix local issue with prefix",
			s:       "prefix Fix #123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "fix local issue with suffix",
			s:       "Fix #123 suffix",
			wantErr: issue.ErrNoIssue,
		},

		{
			name: "remote issue",
			s:    "Issue owner/repo#123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "remote issue with prefix",
			s:       "prefix Issue owner/repo#123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "remote issue with suffix",
			s:       "Issue owner/repo#123 suffix",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "close remote issue",
			s:    "Close owner/repo#123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "close remote issue with prefix",
			s:       "prefix Close owner/repo#123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "close remote issue with suffix",
			s:       "Close owner/repo#123 suffix",
			wantErr: issue.ErrNoIssue,
		},
		{
			name: "fix remote issue",
			s:    "Fix owner/repo#123",
		},
		{
			// can't mix the issue with other text on the line
			name:    "fix remote issue with prefix",
			s:       "prefix Fix owner/repo#123",
			wantErr: issue.ErrNoIssue,
		},
		{
			// can't mix the issue with other text on the line
			name:    "fix remote issue with suffix",
			s:       "Fix owner/repo#123 suffix",
			wantErr: issue.ErrNoIssue,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := issue.Check([]byte(tc.s))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Check() = %v; want %v", err, tc.wantErr)
				t.Logf("data:\n%s", tc.s)
			}
		})
	}
}
