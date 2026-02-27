// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package report_test

import (
	"bytes"
	"flag"
	"testing"

	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/job"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/report"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/scheduler"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/sim"
	"github.com/skhal/lab/book/ostep/cpu/cmd/scheduler/internal/trace"
	"github.com/skhal/lab/go/tests"
)

var update = flag.Bool("update", false, "update golden files")

func TestGenerate_fifo(t *testing.T) {
	tests := []struct {
		name   string
		jobs   []job.Job
		golden tests.GoldenFile
	}{
		{
			name:   "no jobs",
			golden: "testdata/fifo/no_jobs.txt",
		},
		{
			name: "one job arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
			},
			golden: "testdata/fifo/one_job_arrive_immediately.txt",
		},
		{
			name: "one job arrive late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Arrival: 2, Duration: 5}},
			},
			golden: "testdata/fifo/one_job_arrive_late.txt",
		},
		{
			name: "two jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Duration: 3}},
			},
			golden: "testdata/fifo/two_jobs_arrive_immediately.txt",
		},
		{
			name: "two jobs one arrives late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Arrival: 1, Duration: 3}},
			},
			golden: "testdata/fifo/two_jobs_one_arrives_late.txt",
		},
		{
			name: "canon book three jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 10}},
				{ID: 2, Spec: job.Spec{Duration: 10}},
				{ID: 3, Spec: job.Spec{Duration: 10}},
			},
			golden: "testdata/fifo/canon_three_jobs_arrive_immediately.txt",
		},
		{
			name: "canon book three skewed jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 100}},
				{ID: 2, Spec: job.Spec{Duration: 10}},
				{ID: 3, Spec: job.Spec{Duration: 10}},
			},
			golden: "testdata/fifo/canon_three_skewed_jobs_arrive_immediately.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			simu := sim.New(tc.jobs, scheduler.NewFIFO())
			data := report.Data{
				Policy: "fifo",
				Jobs:   tc.jobs,
				Sim:    simu,
				Tracer: trace.NewTracer(simu),
			}

			err := report.Generate(&b, data)

			if err != nil {
				t.Logf("data:\n%v", data)
				t.Fatalf("Generate() unexpected error %v", err)
			}
			got := b.String()
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerate_sjf(t *testing.T) {
	tests := []struct {
		name   string
		jobs   []job.Job
		golden tests.GoldenFile
	}{
		{
			name:   "no jobs",
			golden: "testdata/sjf/no_jobs.txt",
		},
		{
			name: "one job arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
			},
			golden: "testdata/sjf/one_job_arrive_immediately.txt",
		},
		{
			name: "one job arrive late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Arrival: 2, Duration: 5}},
			},
			golden: "testdata/sjf/one_job_arrive_late.txt",
		},
		{
			name: "two jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Duration: 3}},
			},
			golden: "testdata/sjf/two_jobs_arrive_immediately.txt",
		},
		{
			name: "two jobs one arrives late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Arrival: 1, Duration: 3}},
			},
			golden: "testdata/sjf/two_jobs_one_arrives_late.txt",
		},
		{
			name: "canon book three skewed jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 100}},
				{ID: 2, Spec: job.Spec{Duration: 10}},
				{ID: 3, Spec: job.Spec{Duration: 10}},
			},
			golden: "testdata/sjf/canon_three_skewed_jobs_arrive_immediately.txt",
		},
		{
			name: "canon book three skewed jobs arrive late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 100}},
				{ID: 2, Spec: job.Spec{Arrival: 10, Duration: 10}},
				{ID: 3, Spec: job.Spec{Arrival: 10, Duration: 10}},
			},
			golden: "testdata/sjf/canon_three_skewed_jobs_arrive_late.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			simu := sim.New(tc.jobs, scheduler.NewSJF())
			data := report.Data{
				Policy: "sjf",
				Jobs:   tc.jobs,
				Sim:    simu,
				Tracer: trace.NewTracer(simu),
			}

			err := report.Generate(&b, data)

			if err != nil {
				t.Logf("data:\n%v", data)
				t.Fatalf("Generate() unexpected error %v", err)
			}
			got := b.String()
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerate_stcf(t *testing.T) {
	tests := []struct {
		name   string
		jobs   []job.Job
		golden tests.GoldenFile
	}{
		{
			name:   "no jobs",
			golden: "testdata/stcf/no_jobs.txt",
		},
		{
			name: "one job arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
			},
			golden: "testdata/stcf/one_job_arrive_immediately.txt",
		},
		{
			name: "one job arrive late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Arrival: 2, Duration: 5}},
			},
			golden: "testdata/stcf/one_job_arrive_late.txt",
		},
		{
			name: "two jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Duration: 3}},
			},
			golden: "testdata/stcf/two_jobs_arrive_immediately.txt",
		},
		{
			name: "two jobs one arrives late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 5}},
				{ID: 2, Spec: job.Spec{Arrival: 1, Duration: 3}},
			},
			golden: "testdata/stcf/two_jobs_one_arrives_late.txt",
		},
		{
			name: "canon book three skewed jobs arrive late",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 100}},
				{ID: 2, Spec: job.Spec{Arrival: 10, Duration: 10}},
				{ID: 3, Spec: job.Spec{Arrival: 10, Duration: 10}},
			},
			golden: "testdata/stcf/canon_three_skewed_jobs_arrive_late.txt",
		},
		{
			name: "canon book three jobs arrive immediately",
			jobs: []job.Job{
				{ID: 1, Spec: job.Spec{Duration: 10}},
				{ID: 2, Spec: job.Spec{Duration: 10}},
				{ID: 3, Spec: job.Spec{Duration: 10}},
			},
			golden: "testdata/stcf/canon_three_jobs_arrive_immediately.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			simu := sim.New(tc.jobs, scheduler.NewSTCF())
			data := report.Data{
				Policy: "stcf",
				Jobs:   tc.jobs,
				Sim:    simu,
				Tracer: trace.NewTracer(simu),
			}

			err := report.Generate(&b, data)

			if err != nil {
				t.Logf("data:\n%v", data)
				t.Fatalf("Generate() unexpected error %v", err)
			}
			got := b.String()
			if *update {
				tc.golden.Write(t, got)
			}
			if diff := tc.golden.Diff(t, got); diff != "" {
				t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
