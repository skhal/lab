// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin_test

import (
	"testing"

	"github.com/skhal/lab/x/fin/internal/fin"
)

func TestCents_String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		c    fin.Cents
		want string
	}{
		{
			name: "empty",
			want: "0.00",
		},
		{
			name: "no integer part",
			c:    57,
			want: "0.57",
		},
		{
			name: "no fractional part",
			c:    12300,
			want: "123.00",
		},
		{
			name: "mixed parts",
			c:    12345,
			want: "123.45",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.c.String()

			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
