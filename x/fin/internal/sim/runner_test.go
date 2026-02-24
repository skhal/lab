// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/sim"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCycles(t *testing.T) {
	newRecord := func(t *testing.T, year int32, month time.Month) *pb.Record {
		t.Helper()
		m := int32(month)
		return pb.Record_builder{
			Date: pb.Date_builder{
				Year:  &year,
				Month: &m,
			}.Build(),
		}.Build()
	}
	tests := []struct {
		name   string
		data   []*pb.Record
		cycles int
		size   int
		want   [][]*pb.Record
	}{
		{
			name: "insufficient records",
			data: []*pb.Record{
				newRecord(t, 2006, time.January),
			},
			cycles: 1,
			size:   2,
		},
		{
			name: "one cycle size 1",
			data: []*pb.Record{
				newRecord(t, 2006, time.January),
			},
			cycles: 1,
			size:   1,
			want: [][]*pb.Record{
				{
					newRecord(t, 2006, time.January),
				},
			},
		},
		{
			name: "one cycle size 2",
			data: []*pb.Record{
				newRecord(t, 2006, time.January),
				newRecord(t, 2006, time.February),
			},
			cycles: 1,
			size:   2,
			want: [][]*pb.Record{
				{
					newRecord(t, 2006, time.January),
					newRecord(t, 2006, time.February),
				},
			},
		},
		{
			name: "two cycles size 1",
			data: []*pb.Record{
				newRecord(t, 2006, time.January),
				newRecord(t, 2006, time.February),
				newRecord(t, 2006, time.March),
				newRecord(t, 2006, time.April),
				newRecord(t, 2006, time.May),
				newRecord(t, 2006, time.June),
				newRecord(t, 2006, time.July),
				newRecord(t, 2006, time.August),
				newRecord(t, 2006, time.September),
				newRecord(t, 2006, time.October),
				newRecord(t, 2006, time.November),
				newRecord(t, 2006, time.December),
				newRecord(t, 2007, time.January),
			},
			cycles: 2,
			size:   1,
			want: [][]*pb.Record{
				{
					newRecord(t, 2007, time.January),
				},
				{
					newRecord(t, 2006, time.January),
				},
			},
		},
		{
			name: "two cycles size 2",
			data: []*pb.Record{
				newRecord(t, 2006, time.January),
				newRecord(t, 2006, time.February),
				newRecord(t, 2006, time.March),
				newRecord(t, 2006, time.April),
				newRecord(t, 2006, time.May),
				newRecord(t, 2006, time.June),
				newRecord(t, 2006, time.July),
				newRecord(t, 2006, time.August),
				newRecord(t, 2006, time.September),
				newRecord(t, 2006, time.October),
				newRecord(t, 2006, time.November),
				newRecord(t, 2006, time.December),
				newRecord(t, 2007, time.January),
				newRecord(t, 2007, time.February),
			},
			cycles: 2,
			size:   2,
			want: [][]*pb.Record{
				{
					newRecord(t, 2007, time.January),
					newRecord(t, 2007, time.February),
				},
				{
					newRecord(t, 2006, time.January),
					newRecord(t, 2006, time.February),
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(sim.Cycles(tc.data, tc.cycles, tc.size))

			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Collect() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
