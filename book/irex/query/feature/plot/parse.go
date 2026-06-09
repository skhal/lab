// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skhal/lab/book/irex/pb"
)

var (
	// ErrNoSymbol means the plot command does not have a symbol parameter.
	ErrNoSymbol = errors.New("missing symbol")

	// ErrSinceDate means since date has unsupported format.
	ErrSinceDate = errors.New("invalid since date")

	// ErrUntilDate means until date has unsupported format.
	ErrUntilDate = errors.New("invalid until date")
)

var indexByName = map[string]pb.Symbol_Index{
	"spx": pb.Symbol_IDX_SPX,
}

// Parse parses plot-command parameters. It returns the intent upon successful
// parse or error if the parse fails.
func Parse(params string) (*pb.PlotIntent, error) {
	p := &parser{}
	if err := p.Parse(params); err != nil {
		return nil, err
	}
	msg := pb.PlotIntent_builder{
		Symbol: p.Symbol(),
		Since:  p.Since(),
		Until:  p.Until(),
	}.Build()
	return msg, nil
}

type parser struct {
	symbol *pb.Symbol
	since  *pb.Date
	until  *pb.Date
}

// Parse extracts symbol and dates range from the plot-command parameters. The
// dates range is defined by since and until dates, both in YYYY-MM format.
// The parser does not care about the order of the parameters and takes
// multiple trigger words for since and until dates:
// - since: after, from, since
// - until: before, to, until
func (p *parser) Parse(params string) error {
	fields := strings.Fields(strings.ToLower(params))
	if len(fields) == 0 {
		return ErrNoSymbol
	}
	if err := p.parse(fields); err != nil {
		return err
	}
	if p.symbol == nil {
		return ErrNoSymbol
	}
	return nil
}

func (p *parser) parse(fields []string) error {
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		switch {
		case sinceWord[f] && i+1 < len(fields):
			i++ // skip trigger word
			if err := p.parseSince(fields[i]); err != nil {
				return err
			}
		case untilWord[f] && i+1 < len(fields):
			i++ // skip trigger word
			if err := p.parseUntil(fields[i]); err != nil {
				return err
			}
		default:
			idx, ok := indexByName[f]
			if !ok {
				break
			}
			p.symbol = pb.Symbol_builder{Index: &idx}.Build()
		}
	}
	return nil
}

var sinceWord = map[string]bool{
	// keep-sorted start
	"after": true,
	"from":  true,
	"since": true,
	// keep-sorted end
}

var untilWord = map[string]bool{
	// keep-sorted start
	"before": true,
	"to":     true,
	"until":  true,
	// keep-sorted end
}

func (p *parser) parseSince(date string) error {
	d, err := parseDate(date)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrSinceDate, err)
	}
	p.since = d
	return nil
}

func (p *parser) parseUntil(date string) error {
	d, err := parseDate(date)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUntilDate, err)
	}
	p.until = d
	return nil
}

const yearMonthOnly = "2006-01"

func parseDate(s string) (*pb.Date, error) {
	t, err := time.Parse(yearMonthOnly, s)
	if err != nil {
		return nil, err
	}

	var hh, mm, ss, ns int
	t = time.Date(t.Year(), t.Month(), 1, hh, mm, ss, ns, time.UTC)

	years, months, days := 0, 1, 0
	t = t.AddDate(years, months, days)

	years, months, days = 0, 0, -1
	t = t.AddDate(years, months, days)

	return newDate(t), nil
}

func newDate(t time.Time) *pb.Date {
	return pb.Date_builder{
		Year:  new(int32(t.Year())),
		Month: new(int32(t.Month())),
		Day:   new(int32(t.Day())),
	}.Build()
}

// Symbol returns parsed symbol else nil.
func (p *parser) Symbol() *pb.Symbol {
	return p.symbol
}

// Since returns parsed since date else nil.
func (p *parser) Since() *pb.Date {
	return p.since
}

// Until returns parsed until date else nil.
func (p *parser) Until() *pb.Date {
	return p.until
}
