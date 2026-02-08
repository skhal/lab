// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/x/fin/internal/pb"
)

const (
	colDate        = 0
	colSPComposite = 1
	colDividend    = 2
	colEarnings    = 3
)

func parseRow(row []string) (*pb.Record, error) {
	d, err := ParseDate(row[colDate])
	if err != nil {
		return nil, err
	}
	q, err := parseQuote(row)
	if err != nil {
		return nil, err
	}
	r := &pb.Record{}
	r.SetDate(d)
	r.SetQuote(q)
	return r, nil
}

const (
	idxYear  = 0
	idxMonth = 1
)

const (
	january  = 1
	december = 12
)

// ParseDate parses a value from the date cell. It expects the date to be in
// the YYYY.MM form with the exception for October, which uses YYYY.1 format
// (presumably Excel drops zero and treats YYYY.MM as s number with up to 2
// precision digits... so 1234.10 becomes 1234.1).
func ParseDate(s string) (d *pb.Date, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse %s: invalid date, %s", s, err)
		}
	}()
	tokens := strings.SplitN(s, ".", 2)
	if len(tokens) != 2 {
		return nil, errors.New("want ####.##")
	}
	if len(tokens[idxYear]) != 4 {
		return nil, errors.New("want #### for year")
	}
	if n := len(tokens[idxMonth]); n != 1 && n != 2 {
		return nil, errors.New("want # or ## for month")
	}
	y, err := strconv.Atoi(tokens[idxYear])
	if err != nil {
		return nil, err
	}
	if y < 0 {
		return nil, errors.New("negative year")
	}
	// The date format is YYYY.MM except October because YYYY.10 = YYYY.1
	// and YYYY.01 is January.
	var m = 10
	if tokens[1] != "1" {
		m, err = strconv.Atoi(tokens[idxMonth])
		if err != nil {
			return nil, err
		}
		if m < january || m > december {
			return nil, errors.New("invalid month")
		}
	}
	d = new(pb.Date)
	d.SetYear(int32(y))
	d.SetMonth(int32(m))
	return
}

func parseQuote(row []string) (*pb.Quote, error) {
	q := &pb.Quote{}
	for _, item := range []struct {
		c cell
		p centsParser
	}{
		{cell{"sp somposite", colSPComposite}, centsParser{q.SetSpComposite}},
		{cell{"dividend", colDividend}, centsParser{q.SetDividend}},
		{cell{"earnings", colEarnings}, centsParser{q.SetEarnings}},
	} {
		if err := item.p.Parse(row, item.c); err != nil {
			return nil, err
		}
	}
	return q, nil
}

type cell struct {
	name string
	col  int
}

type centsParser struct {
	callback func(*pb.Cents)
}

// Parse parses the cell value as cents. It passes the value to the callback
// upon successful parse or returns an error.
func (p centsParser) Parse(row []string, cell cell) error {
	s := row[cell.col]
	c, err := ParseCents(s)
	if err != nil {
		err = fmt.Errorf("parse %s: invalid %s, %s", cell.name, s, err)
		return err
	}
	p.callback(c)
	return nil
}

const (
	idxIntegerPart     = 0
	idxFranctionalPart = 1
)

// ParseCents parses the string ####.## as cents, i.e. the integral part is
// multiplied by 100 and summed with fractional part (i*100 + f).
func ParseCents(s string) (c *pb.Cents, err error) {
	tokens := strings.SplitN(s, ".", 2)
	if len(tokens) != 2 {
		return nil, errors.New("want ###.##")
	}
	// Unlike the date field with October (YYYY.10 = YYYY.1), the prices do
	// add zero, i.e., ####.10 != ####.1, it still keeps 10 for 10 cents.
	if len(tokens[idxFranctionalPart]) != 2 {
		return nil, errors.New("want ## for fractional part")
	}
	n, err := strconv.Atoi(tokens[idxIntegerPart])
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, errors.New("negative integer part")
	}
	m, err := strconv.Atoi(tokens[idxFranctionalPart])
	if err != nil {
		return nil, err
	}
	if m < 0 {
		return nil, errors.New("negative fractional part")
	}
	c = &pb.Cents{}
	c.SetCents(int32(n*100 + m))
	return c, nil
}
