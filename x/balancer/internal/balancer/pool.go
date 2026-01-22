// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"bytes"
	"fmt"
	"math"
	"sync"
)

type RoundRobinPool struct {
	workers []*Worker
	wg      *sync.WaitGroup

	idx   int
	stats []int
}

func newRoundRobinPool(numWorkers int) *RoundRobinPool {
	return &RoundRobinPool{
		workers: make([]*Worker, 0, numWorkers),
		wg:      new(sync.WaitGroup),
		stats:   make([]int, numWorkers),
	}
}

func (p *RoundRobinPool) Start() {
	for len(p.workers) < cap(p.workers) {
		w := p.startWorker()
		p.workers = append(p.workers, w)
	}
}

func (p *RoundRobinPool) startWorker() *Worker {
	w := newWorker()
	p.wg.Go(w.Run)
	return w
}

func (p *RoundRobinPool) Dispatch(req *Request) {
	w := p.selectWorker()
	w.Send(req)
}

func (p *RoundRobinPool) selectWorker() *Worker {
	w := p.workers[p.idx]
	p.stats[p.idx] += 1
	p.idx += 1
	p.idx %= len(p.workers)
	return w
}

func (p *RoundRobinPool) Stop() {
	for _, w := range p.workers {
		w.Stop()
	}
}

func (p *RoundRobinPool) Wait() {
	p.wg.Wait()
}

type PoolStats struct {
	load []int
}

func newPoolStats(load []int) *PoolStats {
	return &PoolStats{
		load: load,
	}
}

func (ps *PoolStats) String() string {
	mean, std := calculateStatistics(ps.load)
	var buf bytes.Buffer
	for id, load := range ps.load {
		if id > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", load)
	}
	fmt.Fprintf(&buf, " %.2f %.2f", mean, std)
	return buf.String()
}

func (p *RoundRobinPool) Stats() *PoolStats {
	load := make([]int, len(p.stats))
	copy(load, p.stats)
	return newPoolStats(load)
}

func calculateStatistics(nn []int) (mean float64, std float64) {
	switch len(nn) {
	case 0:
		return 0, 0
	case 1:
		return float64(nn[0]), 0
	}
	var (
		sum   int64
		sumsq int64
	)
	for _, n := range nn {
		n := int64(n)
		sum += n
		sumsq += n * n
	}
	sz := int64(len(nn))
	mean = float64(sum) / float64(sz)
	std = math.Sqrt(float64(sz*sumsq-sum*sum) / float64(sz*(sz-1)))
	return
}
