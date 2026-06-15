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
	// ErrNoSymbol means the query does not have a symbol parameter, e.g.:
	//		plot
	ErrNoSymbol = errors.New("missing symbol")

	// ErrMultipleSymbol means the query has multiple symbols, e.g.:
	//		plot idx idx
	ErrMultipleSymbol = errors.New("multiple symbols")

	// ErrMultipleIndexMetric means the query has multiple index metrics, e.g.:
	// 		plot idx div earn
	ErrMultipleIndexMetric = errors.New("multiple symbol metrics")

	// ErrNotIndex means the symbol is not an index when parsing an index metric,
	// e.g.:
	//		plot cpi div
	ErrNotIndex = errors.New("not index")

	// ErrSinceDate means since date has unsupported format, e.g.:
	//		plot spx since 98-01
	ErrSinceDate = errors.New("invalid since date")

	// ErrUntilDate means until date has unsupported format, e.g.:
	//		plot spx until 98-01
	ErrUntilDate = errors.New("invalid until date")
)

var indexIdentifierByName = map[string]pb.Symbol_Index_ID{
	"spx": pb.Symbol_Index_ID_SPX,
}

var indexMetricByName = map[string]pb.Symbol_Index_Metric{
	// keep-sorted start
	"div":  pb.Symbol_Index_MET_DIV,
	"earn": pb.Symbol_Index_MET_EARN,
	// keep-sorted end
}

var marketMetricByName = map[string]pb.Symbol_Market_Metric{
	"cpi": pb.Symbol_Market_MET_CPI,
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
			if err := p.parseSymbol(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) parseSymbol(s string) error {
	if idx, ok := indexIdentifierByName[s]; ok {
		return p.setIndex(idx)
	}

	if m, ok := indexMetricByName[s]; ok {
		return p.setIndexMetric(m)
	}

	if m, ok := marketMetricByName[s]; ok {
		return p.setMarketMetric(m)
	}

	return nil
}

func (p *parser) setIndex(idx pb.Symbol_Index_ID) error {
	if p.symbol != nil {
		return fmt.Errorf("%w: %s - new %s", ErrMultipleSymbol, p.symbol, idx)
	}
	p.symbol = pb.Symbol_builder{
		Index: pb.Symbol_Index_builder{
			Id: &idx,
		}.Build(),
	}.Build()
	return nil
}

func (p *parser) setIndexMetric(m pb.Symbol_Index_Metric) error {
	if p.symbol == nil {
		return fmt.Errorf("%w: index metric %s", ErrNoSymbol, m)
	}
	if !p.symbol.HasIndex() {
		return fmt.Errorf("%w: %s - index metric %s", ErrNotIndex, p.symbol, m)
	}
	if p.symbol.GetIndex().HasMetric() {
		return fmt.Errorf("%w: %s - new metric %s", ErrMultipleIndexMetric, p.symbol, m)
	}
	p.symbol.GetIndex().SetMetric(m)
	return nil
}

func (p *parser) setMarketMetric(m pb.Symbol_Market_Metric) error {
	if p.symbol != nil {
		return fmt.Errorf("%w: %s - market metric %s", ErrMultipleSymbol, p.symbol, m)
	}
	p.symbol = pb.Symbol_builder{
		Market: pb.Symbol_Market_builder{
			Metric: &m,
		}.Build(),
	}.Build()
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
