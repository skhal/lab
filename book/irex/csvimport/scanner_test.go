// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport_test

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/csvimport"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name      string
		csv       string
		skipLines int
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
1990.01,1.11
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 111),
			},
		},
		{
			name:      "one record skip lines",
			skipLines: 1,
			csv: `
1990.01,1.11
1990.02,2.22
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.February, 28, 222),
			},
		},
		{
			name: "two records",
			csv: `
1990.01,1.11
1990.02,2.22
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 111),
				newQuote(t, 1990, time.February, 28, 222),
			},
		},
		{
			name:      "negative skip lines",
			skipLines: -1,
			wantErr:   csvimport.ErrScan,
		},
		{
			name:      "skip invalid lines",
			skipLines: 2,
			csv: `
foo,bar
baz
1990.01,1.11
`,
			wantErr: csvimport.ErrScan,
		},
		{
			name: "invalid lines",
			csv: `
1990.01,1.11
1990.02
`,
			want: []*pb.Quote{
				newQuote(t, 1990, time.January, 31, 111),
			},
			wantErr: csvimport.ErrScan,
		},
		{
			name: "invalid date",
			csv: `
1000.01,1.11
`,
			wantErr: csvimport.ErrDate,
		},
		{
			name: "invalid spx",
			csv: `
1990.01,1.ab
`,
			wantErr: csvimport.ErrSPX,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rstr := strings.NewReader(tc.csv)
			rcsv := csv.NewReader(rstr)
			sc := csvimport.NewScanner(rcsv)
			sc.SkipLines = tc.skipLines

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

func ExampleScanner() {
	data := `
1990.01,1.01
1990.02,1.02
`
	r := csv.NewReader(strings.NewReader(data))
	sc := csvimport.NewScanner(r)
	var quotes []*pb.Quote
	for sc.Next() {
		quotes = append(quotes, sc.Quote())
	}
	if err := sc.Err(); err != nil {
		fmt.Println(err)
		return
	}
	for _, q := range quotes {
		fmt.Printf("%s\n", q)
	}
	// Output:
	// date:{year:1990  month:1  day:31}  spx:{value:101}
	// date:{year:1990  month:2  day:28}  spx:{value:102}
}
