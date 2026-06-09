// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package market

import (
	"iter"
	"slices"

	"github.com/skhal/lab/book/irex/pb"
)

// DateRange filters quotes to the dates starting from the since date and
// before the until date (excluded). A missing boundary has no effect on the
// restricting, i.e. a nil since-date makes the DateRange to pass all quotes
// until the until-date.
type DateRange struct {
	// Since restricts quotes to the dates equal or after this date.
	Since *pb.Date

	// Until restricts quotes to the date before this date.
	Until *pb.Date
}

// Quotes returns an iterator over quotes whose dates are within the DatRange
// since- and until-dates.
func (dr *DateRange) Quotes(quotes []*pb.Quote) iter.Seq[*pb.Quote] {
	return func(yield func(*pb.Quote) bool) {
		quotes = skipUntil(dr.Since, quotes)
		for q := range rangeUntil(dr.Until, quotes) {
			if !yield(q) {
				break
			}
		}
	}
}

func skipUntil(until *pb.Date, quotes []*pb.Quote) []*pb.Quote {
	if until == nil {
		return quotes
	}
	for len(quotes) != 0 {
		if cmp := compareDate(quotes[0].GetDate(), until); cmp >= 0 {
			break
		}
		quotes = quotes[1:]
	}
	return quotes
}

func rangeUntil(until *pb.Date, quotes []*pb.Quote) iter.Seq[*pb.Quote] {
	if until == nil {
		return slices.Values(quotes)
	}
	return func(yield func(*pb.Quote) bool) {
		for _, q := range quotes {
			if cmp := compareDate(q.GetDate(), until); cmp >= 0 {
				break
			}
			if !yield(q) {
				break
			}
		}
	}
}

// compareDate returns a negative number if a should be before b, a positive
// number if a should be after b, else zero.
func compareDate(a, b *pb.Date) int {
	if d := int(a.GetYear() - b.GetYear()); d != 0 {
		return d
	}
	if d := int(a.GetMonth() - b.GetMonth()); d != 0 {
		return d
	}
	return int(a.GetDay() - b.GetDay())
}
