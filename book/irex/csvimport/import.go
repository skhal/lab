// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport

import (
	"encoding/csv"
	"errors"
	"fmt"

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
	var quotes []*pb.Quote
	sc := scanner.New(r)
	sc.SkipLines = imp.skipLines
	sc.ScanLines = imp.scanLines
	for sc.Next() {
		quotes = append(quotes, sc.Quote())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrImport, err)
	}
	return quotes, nil
}
