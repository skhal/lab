// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"bytes"
	"fmt"
	"io"

	"github.com/skhal/lab/x/fin/internal/fin"
)

// StrategyInfo summarizes the results of a strategy. It includes strategy name,
// description, and starting/end quotes.
type StrategyInfo struct {
	Name        string    // strategy name
	Description string    // strategy description
	Start       fin.Quote // strategy starting quote
	End         fin.Quote // strategy result
}

// Strategy generates a report for a single strategy.
func Strategy(w io.Writer, info StrategyInfo) error {
	return tmpls.ExecuteTemplate(w, "strategy.txt", info)
}

// Strategies lists per-strategy reports.
func Strategies(w io.Writer, infos []*StrategyInfo) error {
	if err := tmpls.ExecuteTemplate(w, "strategies.txt", infos); err != nil {
		return err
	}
	return strategiesPerformance(w, infos)
}

func strategiesPerformance(w io.Writer, infos []*StrategyInfo) error {
	header := func() {
		var b bytes.Buffer
		fmt.Fprint(&b, "    ")
		for i := range infos {
			if i > 0 {
				fmt.Fprint(&b, " | ")
			}
			fmt.Fprintf(&b, " [%d]", i)
		}
		fmt.Fprintln(&b)
		io.Copy(w, &b)
	}
	row := func(r int, rinfo *StrategyInfo) {
		var b bytes.Buffer
		fmt.Fprintf(&b, "[%d] ", r)
		for c, cinfo := range infos {
			if c > 0 {
				fmt.Fprint(&b, " | ")
			}
			if c != r {
				rate := float64(cinfo.End.Balance) / float64(rinfo.End.Balance)
				fmt.Fprintf(&b, "%.2f", rate)
			} else {
				fmt.Fprintf(&b, "%4s", "")
			}
		}
		fmt.Fprintln(&b)
		io.Copy(w, &b)
	}
	list := func() {
		for i, info := range infos {
			fmt.Fprintf(w, "[%d] %s\n", i, info.Name)
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Performance:")
	for r, rinfo := range infos {
		if r == 0 {
			header()
		}
		row(r, rinfo)
	}
	fmt.Fprintln(w)
	list()
	return nil
}
