// Copyright 2025 Samvel Khalatyan. All rights reserved.

package stripzero_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/9/stripzero"
)

func TestCleanOptimal(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m := DeepCopy(t, tc.m)

			stripzero.CleanOptimal(m)

			if diff := cmp.Diff(tc.want, m); diff != "" {
				t.Errorf("stripzero.CleanOptimal(...) mismatch (-want, +got):\n%s", diff)
				t.Logf("Input:\n%s", tc.m)
			}
		})
	}
}
