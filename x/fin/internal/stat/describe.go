// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math/big"

	"github.com/skhal/lab/x/fin/internal/fin"
)

// Description holds descriptive statistics to summarize a sample.
type Description struct {
	// Max is the maximum value of the sample.
	Max fin.Cents
	// Min is the minimum value of the sample
	Min fin.Cents

	// Avg is the arithmetic average of the sample.
	Avg fin.Cents
	// Med is the median of the sample.
	Med fin.Cents

	// Std is the standard deviation of the sample.
	Std fin.Cents
}

// Describe generates a descriptive statistics.
func Describe(cc []fin.Cents) Description {
	if len(cc) == 0 {
		return Description{}
	}
	d := &describer{
		sum:   new(big.Int),
		sumsq: new(big.Int),
	}
	for _, c := range cc {
		d.Add(c)
	}
	return Description{
		Max: d.Max(),
		Min: d.Min(),
		Avg: d.Avg(),
		Med: median(cc),
		Std: d.Std(),
	}
}

func median(cc []fin.Cents) fin.Cents {
	switch n := len(cc); {
	case n == 1:
		return cc[0]
	case n%2 == 0:
		return (cc[n/2-1] + cc[n/2]) / 2
	case n%2 == 1:
		return cc[n/2]
	}
	return 0
}

type describer struct {
	min fin.Cents
	max fin.Cents

	n     int
	sum   *big.Int
	sumsq *big.Int

	average  *fin.Cents
	variance *fin.Cents
	stddev   *fin.Cents
}

// Add appends a datapoint to the sample.
func (d *describer) Add(c fin.Cents) {
	if d.n == 0 {
		d.min = c
		d.max = c
	} else {
		if c < d.min {
			d.min = c
		}
		if c > d.max {
			d.max = c
		}
	}
	d.n += 1
	bc := big.NewInt(int64(c))
	d.sum.Add(d.sum, bc)
	sq := big.NewInt(int64(c))
	sq.Mul(sq, bc)
	d.sumsq.Add(d.sumsq, sq)
}

// Avg calculates average of the sample.
func (d *describer) Avg() fin.Cents {
	if d.average == nil {
		d.updateAverage()
	}
	return *d.average
}

func (d *describer) updateAverage() {
	var avg big.Int
	avg.Div(d.sum, big.NewInt(int64(d.n)))
	x := fin.Cents(avg.Int64())
	d.average = &x
}

// Max returns maximum value of the sample.
func (d *describer) Max() fin.Cents {
	return d.max
}

// Min returns minimum value of the sample.
func (d *describer) Min() fin.Cents {
	return d.min
}

// Std calculates standard deviation of the sample.
func (d *describer) Std() fin.Cents {
	if d.stddev == nil {
		d.updateStandardDeviation()
	}
	return *d.stddev
}

func (d *describer) updateStandardDeviation() {
	v := d.calculateVariance()
	v.Sqrt(v)
	x := fin.Cents(v.Int64())
	d.stddev = &x
}

func (d *describer) calculateVariance() *big.Int {
	if d.n < 2 {
		return big.NewInt(0)
	}
	lv := big.NewInt(int64(d.n))
	lv.Mul(lv, d.sumsq)
	sumsq := big.NewInt(0)
	sumsq.Add(sumsq, d.sum)
	sumsq.Mul(sumsq, d.sum)
	lv.Sub(lv, sumsq)
	lv.Div(lv, big.NewInt(int64(d.n*(d.n-1))))
	return lv
}
