// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report

import (
	"io"

	"github.com/skhal/lab/x/fin/internal/sim"
)

// SimulationData is input to generate a simulation report.
type SimulationData struct {
	// Name is the strategy name.
	Name string
	// Desc is the strategy description.
	Desc string
	// Result is descriptive statistics, summarizing strategy results.
	Result *sim.Result
}

// Simulation generates a report on SimulationData.
func Simulation(w io.Writer, data []*SimulationData) error {
	return tmpls.ExecuteTemplate(w, "simulation.txt", data)
}
