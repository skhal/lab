// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport_test

import (
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestImport(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		want    []*pb.Quote
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "no data",
			csv: `
1,header
2,header
3,header
4,header
5,header
6,header
7,header
8,header
`,
		},
		{
			name: "non-empty data",
			csv: `
1-date,spx,dividend
2-date,spx,dividend
3-date,spx,dividend
4-date,spx,dividend
5-date,spx,dividend
6-date,spx,dividend
7-date,spx,dividend
8-date,spx,dividend
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
				newQuote(t, 1990, time.February, 28, 201, 202),
			},
		},
		{
			name: "data error",
			csv: `
1-date,spx,dividend
2-date,spx,dividend
3-date,spx,dividend
4-date,spx,dividend
5-date,spx,dividend
6-date,spx,dividend
7-date,spx,dividend
8-date,spx,dividend
1990.01,1.01,1.02
1990.02,abc,2.02
`,
			wantErr: csvimport.ErrImport,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := csv.NewReader(strings.NewReader(tc.csv))

			got, err := csvimport.Import(r)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v', want '%v'", err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Logf("csv:\n%s", tc.csv)
			}
		})
	}
}

func newQuote(t *testing.T, year int32, month time.Month, day int32, spx, div int32) *pb.Quote {
	t.Helper()
	return pb.Quote_builder{
		Date: newDate(t, year, month, day),
		Spx:  pb.Cent_builder{Value: &spx}.Build(),
		Div:  pb.Cent_builder{Value: &div}.Build(),
	}.Build()
}

func newDate(t *testing.T, year int32, month time.Month, day int32) *pb.Date {
	t.Helper()
	return pb.Date_builder{
		Year:  &year,
		Month: new(int32(month)),
		Day:   &day,
	}.Build()
}
