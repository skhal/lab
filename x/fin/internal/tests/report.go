// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"fmt"
	"strings"

	"github.com/skhal/lab/x/fin/internal/report"
)

// InfoStringer makes [report.StrategyInfo] printable for tests.
type InfoStringer report.StrategyInfo

// String implements [fmt.Stringer] interface.
func (info InfoStringer) String() string {
	b := new(strings.Builder)
	fmt.Fprintln(b, "Name:", info.Name)
	fmt.Fprintln(b, "Description:", info.Description)
	fmt.Fprintln(b, "Start:", info.Start)
	fmt.Fprint(b, "End: ", info.End)
	return b.String()
}
