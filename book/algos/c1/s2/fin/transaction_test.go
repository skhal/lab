// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/algos/c1/s2/fin"
)

func ExampleTransaction_formatScan() {
	var s string
	var tx *fin.Transaction
	{
		tx := &fin.Transaction{
			Customer: "example",
			Date:     time.Date(2000, time.March, 10, 0, 0, 0, 0, time.UTC),
			Amount:   1.23,
		}
		s = tx.String()
	}
	fmt.Println(s)
	{
		tx = new(fin.Transaction)
		if _, err := fmt.Sscan(s, tx); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println(tx)
	// Output:
	// example 3/10/2000 1.23
	// example 3/10/2000 1.23
}

func TestTransaction_String(t *testing.T) {
	tests := []struct {
		name string
		tx   *fin.Transaction
		want string
	}{
		{
			name: "empty",
			want: "<nil>",
		},
		{
			name: "non empty",
			tx: &fin.Transaction{
				Customer: "test",
				Date:     newDate(t, 2009, time.November, 18),
				Amount:   1.23,
			},
			want: "test 11/18/2009 1.23",
		},
		{
			name: "amount precision 2",
			tx: &fin.Transaction{
				Customer: "test",
				Date:     newDate(t, 2009, time.November, 18),
				Amount:   1.234,
			},
			want: "test 11/18/2009 1.23",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.tx.String()

			if tc.want != got {
				t.Errorf("fin.Transaction.String() = %q; want %q", got, tc.want)
			}
		})
	}
}

func TestTransaction_Scan(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *fin.Transaction
		wantErr error
	}{
		{
			name:    "empty",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
		{
			name: "non empty",
			s:    "test 11/18/2009 1.23",
			want: &fin.Transaction{
				Customer: "test",
				Date:     newDate(t, 2009, time.November, 18),
				Amount:   1.23,
			},
		},
		{
			name: "mixed spacing",
			s:    "test 	11/18/2009  1.23",
			want: &fin.Transaction{
				Customer: "test",
				Date:     newDate(t, 2009, time.November, 18),
				Amount:   1.23,
			},
		},
		{
			name:    "customer missing",
			s:       "11/18/2009 1.23",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
		{
			name:    "date missing",
			s:       "test 1.23",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
		{
			name:    "date wrong format",
			s:       "test 11-18-2009 1.23",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
		{
			name:    "amount missing",
			s:       "test 11/18/2009",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
		{
			name:    "amount wrong format",
			s:       "test 11/18/2009 abc",
			want:    new(fin.Transaction),
			wantErr: fin.ErrFormat,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := new(fin.Transaction)

			_, err := fmt.Sscanf(tc.s, "%v", got)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("(*fin.Transaction).Scan() = _, %v; want error %v", err, tc.wantErr)
				t.Logf("string: %q", tc.s)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(*fin.Transaction).Scan() mismatch (-want, +got):\n%s", diff)
				t.Logf("string: %q", tc.s)
			}
		})
	}
}

func newDate(t *testing.T, year int, month time.Month, day int) time.Time {
	t.Helper()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
