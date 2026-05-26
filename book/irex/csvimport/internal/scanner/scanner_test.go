// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner_test

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport/internal/scanner"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name      string
		csv       string
		skipLines int
		scanLines int
		want      []*pb.Quote
		wantErr   error
	}{
		{
			name: "empty",
		},
		{
			name:      "empty with skip lines",
			skipLines: 1,
		},
		{
			name: "one record",
			csv: `
1990.01,1.01,1.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
		{
			name:      "one record skip lines",
			skipLines: 1,
			csv: `
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.February, 28, 201, 202),
			},
		},
		{
			name: "two records",
			csv: `
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
				newQuote(t, 1990, time.February, 28, 201, 202),
			},
		},
		{
			name:      "negative skip lines",
			skipLines: -1,
			wantErr:   scanner.ErrScan,
		},
		{
			name:      "skip invalid lines",
			skipLines: 2,
			csv: `
foo,bar
baz
1990.01,1.01,1.02
`,
			wantErr: scanner.ErrScan,
		},
		{
			name: "invalid lines",
			csv: `
1990.01,1.01,1.02
1990.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
			wantErr: scanner.ErrScan,
		},
		{
			name: "invalid date",
			csv: `
1000.01,1.01,1.02
`,
			wantErr: scanner.ErrScan,
		},
		{
			name: "invalid spx",
			csv: `
1990.01,1.ab,1.02
`,
			wantErr: scanner.ErrScan,
		},
		{
			name: "invalid dividend",
			csv: `
1990.01,1.01,1.ab
`,
			wantErr: scanner.ErrScan,
		},
		{
			name: "scan lines one",
			csv: `
1990.01,1.01,1.02
1990.02,2.01,2.02
`,
			scanLines: 1,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 101, 102),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rstr := strings.NewReader(tc.csv)
			rcsv := csv.NewReader(rstr)
			sc := scanner.New(rcsv)
			sc.SkipLines = tc.skipLines
			sc.ScanLines = tc.scanLines

			var got []*pb.Quote
			for sc.Next() {
				got = append(got, sc.Quote())
			}
			err := sc.Err()

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

func TestScanner_Next_noop_on_err(t *testing.T) {
	s := `
1990.01,1.01,1.02
1990.02,ab,2.02
1990.03,3.01,3.02
`
	rcsv := csv.NewReader(strings.NewReader(s))
	sc := scanner.New(rcsv)

	var got []*pb.Quote
	for sc.Next() {
		got = append(got, sc.Quote())
	}
	sc.Next()

	if err, want := sc.Err(), scanner.ErrScan; !errors.Is(err, want) {
		t.Errorf("unexpected error '%v', want '%v'", err, want)
	}
	want := []*pb.Quote{
		newQuote(t, 1990, time.January, 31, 101, 102),
	}
	if d := cmp.Diff(want, got, protocmp.Transform()); d != "" {
		t.Errorf("mismatch (-want +got):\n%s", d)
	}
}

func ExampleScanner() {
	data := `
1990.01,1.01,1.02
1990.02,2.01,2.02
`
	r := csv.NewReader(strings.NewReader(data))
	sc := scanner.New(r)
	var quotes []*pb.Quote
	for sc.Next() {
		quotes = append(quotes, sc.Quote())
	}
	if err := sc.Err(); err != nil {
		fmt.Println(err)
		return
	}
	for _, q := range quotes {
		fmt.Printf("%v\n", q)
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
