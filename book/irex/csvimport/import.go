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

// SkipLines is the number of header lines to skip in CSV data.
const SkipLines = 8

// Import parses CSV data into a list of quotes. It skips first [SkipLines] for
// header lines.
// It returns an error if data parsing fails.
func Import(r *csv.Reader) ([]*pb.Quote, error) {
	var quotes []*pb.Quote
	sc := scanner.New(r)
	sc.SkipLines = SkipLines
	for sc.Next() {
		quotes = append(quotes, sc.Quote())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrImport, err)
	}
	return quotes, nil
}
