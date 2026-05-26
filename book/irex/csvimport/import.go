// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"slices"

	"github.com/skhal/lab/book/irex/csvimport/internal/scanner"
	"github.com/skhal/lab/book/irex/pb"
)

// ErrImport means there was an error in importing CSV data.
var ErrImport = errors.New("import fail")

// Opt configures the importer in some way.
type Opt func(*importer)

// WithSkipLines makes importer skip first n lines.
func WithSkipLines(n int) Opt {
	return func(imp *importer) {
		imp.skipLines = n
	}
}

// WithScanLines makes importer scan at most n lines including skipped lines.
func WithScanLines(n int) Opt {
	return func(imp *importer) {
		imp.scanLines = n
	}
}

// Import parses CSV data into a list of quotes. It returns an error if data
// parsing fails.
func Import(r *csv.Reader, opts ...Opt) ([]*pb.Quote, error) {
	imp := &importer{}
	for _, o := range opts {
		o(imp)
	}
	return imp.Import(r)
}

type importer struct {
	skipLines int
	scanLines int
}

// Import parses CSV data into a list of quotes. It skip first skipLines and
// scans at most scanLines, including the skipLines. A zero value for skipLines
// or scanLines disables skips or limits.
// It returns an error if data parsing fails.
func (imp *importer) Import(r *csv.Reader) ([]*pb.Quote, error) {
	quotes, err := imp.scan(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrImport, err)
	}
	imp.sort(quotes)
	if err := imp.validate(quotes); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrImport, err)
	}
	return quotes, nil
}

func (imp *importer) scan(r *csv.Reader) ([]*pb.Quote, error) {
	var quotes []*pb.Quote
	sc := scanner.New(r)
	sc.SkipLines = imp.skipLines
	sc.ScanLines = imp.scanLines
	for sc.Next() {
		quotes = append(quotes, sc.Quote())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return quotes, nil
}

// sort sorts quotes by dates in ascending order.
func (imp *importer) sort(quotes []*pb.Quote) {
	slices.SortStableFunc(quotes, func(a, b *pb.Quote) int {
		return compareDate(a.GetDate(), b.GetDate())
	})
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

// validate checks that quotes have not duplicate dates and every i-th quote is
// for the next month of (i-1)th quote.
func (imp *importer) validate(quotes []*pb.Quote) error {
	if len(quotes) < 2 {
		return nil
	}
	a := quotes[0]
	for i := 1; i < len(quotes); i++ {
		b := quotes[i]
		if err := validateDateGap(a.GetDate(), b.GetDate()); err != nil {
			return err
		}
		a = b
	}
	return nil
}

// validateDateGap ensures that the date b is in the next month of a,
// regardless of the day. It accounts only for the year and month.
// It returns an error if the condition does not hold.
func validateDateGap(a, b *pb.Date) error {
	switch b.GetYear() {
	case a.GetYear():
		// same year, want b.Month be next of a.Month
		if b.GetMonth()-a.GetMonth() == 1 {
			return nil
		}
	case a.GetYear() + 1:
		// b is the next year wrt a. A should be Dec, b should be Jan
		if a.GetMonth() == 12 && b.GetMonth() == 1 {
			return nil
		}
	}
	return fmt.Errorf("invalid dates %s and %s", (*date)(a), (*date)(b))
}

// date is printable pb.Date.
type date pb.Date

// String prints the date in YYYY.MM format.
func (d *date) String() string {
	return fmt.Sprintf("%d.%d", (*pb.Date)(d).GetYear(), (*pb.Date)(d).GetMonth())
}
