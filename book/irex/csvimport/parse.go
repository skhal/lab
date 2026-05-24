// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvimport

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/skhal/lab/book/irex/pb"
)

var (
	// ErrNoDate means date value is missing in the record.
	ErrNoDate = errors.New("missing date")

	// ErrDate means the date value has invalid format (####.##) or incorrect,
	// e.g. not a number or the year is below [MinYear]. October is an exception
	// to the date format, it is represented by ####.1.
	ErrDate = errors.New("invalid date")

	// ErrYear means the date has invalid year
	ErrYear = errors.New("invalid year")

	// ErrMonth means the date has invalid month.
	ErrMonth = errors.New("invalid month")
)

var (
	// ErrNoSPX means SPX value is missing in the record.
	ErrNoSPX = errors.New("missing SPX")

	// ErrSPX means SPX value has invalid format. It should be ####.## with
	// dollar and cents.
	ErrSPX = errors.New("invalid SPX")
)

var (
	// ErrNoDividend means Dividend value is missing in the record.
	ErrNoDividend = errors.New("missing dividend")

	// ErrDividend means Dividend value has invalid format. See [ErrSPX] for
	// format description.
	ErrDividend = errors.New("invalid dividend")
)

var (
	// ErrFormat means the field has invalid format.
	ErrFormat = errors.New("invalid format")

	// ErrDollar means dollar part of the value in ####.## fails to parse.
	ErrDollar = errors.New("invalid dollar")

	// ErrCent means cent part of the value in ####.## fails to parse.
	ErrCent = errors.New("invalid cent")
)

// MinYear is the minimal supported year. Any date before the first day of this
// year is considered invalid.
const MinYear = 1900

// Parse converts a CSV record (row) into a quote. It returns and error if
// the row is invalid, i.e. it misses fields or the fields are invalid.
func Parse(rec []string) (*pb.Quote, error) {
	date, err := record(rec).Date()
	if err != nil {
		return nil, err
	}
	spx, err := record(rec).SPX()
	if err != nil {
		return nil, err
	}
	div, err := record(rec).Dividend()
	if err != nil {
		return nil, err
	}
	q := pb.Quote_builder{Date: date, Spx: spx, Div: div}.Build()
	return q, nil
}

// record is a CSV row with date, S&P Composite index, etc.
type record []string

const (
	idxDate = iota
	idxSPX
	idxDividend
)

// Date returns the date value of the record.
// It returns an error in case the date is missing or invalid.
func (rec record) Date() (*pb.Date, error) {
	if len(rec) <= idxDate {
		return nil, ErrNoDate
	}
	d, err := ParseDate(rec[idxDate])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDate, err)
	}
	return d, nil
}

// ParseDate parses date string ####.## into year, month, and the last day of
// the month. The year has minimum threshold of [MinYear]. The month is in the
// range (1, 12) inclusive, represented by two digits except October, which is
// given by .1, e.g. 1990.1 is Oct 1990.
// It returns an error if parsing fails.
func ParseDate(s string) (*pb.Date, error) {
	tokens := strings.SplitN(s, ".", 2)
	if len(tokens) != 2 {
		return nil, fmt.Errorf("%w: %s must be YYYY.MM", ErrFormat, s)
	}
	year, err := date(tokens).Year()
	if err != nil {
		return nil, err
	}
	month, err := date(tokens).Month()
	if err != nil {
		return nil, err
	}
	day := lastDayOf(year, month)
	dt := pb.Date_builder{
		Year:  new(int32(year)),
		Month: new(int32(month)),
		Day:   new(int32(day)),
	}.Build()
	return dt, nil
}

func lastDayOf(year, month int) int {
	var (
		day                  = 1
		hour, min, sec, nsec int
	)
	// any timezone will work
	t := time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC)
	t = t.AddDate(0, 1, 0)  // first of the next month
	t = t.AddDate(0, 0, -1) // get last day of the previoud month
	return t.Day()
}

// date is a pair of the year and month components of ####.## value.
type date []string

const (
	idxDateYear = iota
	idxDateMonth
)

// Year returns the year of the date, the value before the dot, e.g. 1990 in
// 1990.11. The year has a minimum year cutoff [MinYear].
// It returns an error if year is invalid.
func (d date) Year() (int, error) {
	n, err := strconv.Atoi(d[idxDateYear])
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrYear, err)
	}
	if n < MinYear {
		return 0, fmt.Errorf("%w: %d below minimum %d", ErrYear, n, MinYear)
	}
	return n, nil
}

// Month returns the month number of the date, the value after dot, e.g. 11 in
// 1990.11. The month lays in the range [1, 12] inclusive. Months take two
// digits after dot, e.g. 1990.03 for March. Keep in mind that October stores
// as a single digit, e.g. 1990.1 because it has number 10 that converts to .1,
// not .10.
// The function returns an error if the month is invalid.
func (d date) Month() (int, error) {
	s := d[idxDateMonth]
	if len(s) == 1 {
		// "1" means month 10 or October
		if s != "1" {
			return 0, fmt.Errorf("%w: %s", ErrMonth, s)
		}
		return 10, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrMonth, err)
	}
	if n < 1 || n > 12 {
		return 0, fmt.Errorf("%w: %s", ErrMonth, s)
	}
	return n, nil
}

// SPX returns the cents value of Standard&Poor Composite index value of the
// record.
// It returns an error if the value is missing or invalid.
func (rec record) SPX() (*pb.Cent, error) {
	if len(rec) <= idxSPX {
		return nil, ErrNoSPX
	}
	c, err := ParseCent(rec[idxSPX])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSPX, err)
	}
	return c, nil
}

// Dividend parses the Dividend value of the record.
// It returns an error if the value is missing or invalid.
func (rec record) Dividend() (*pb.Cent, error) {
	if len(rec) <= idxDividend {
		return nil, ErrNoDividend
	}
	c, err := ParseCent(rec[idxDividend])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDividend, err)
	}
	return c, nil
}

// ParseCent parses string ####.## as cents. It returns an error if the value
// has invalid format such as cents or dollars are missing, the cents must have
// two digits, or the value is not a number.
func ParseCent(s string) (*pb.Cent, error) {
	tokens := strings.SplitN(s, ".", 2)
	if len(tokens) != 2 {
		return nil, fmt.Errorf("%w: %s must be #.##", ErrFormat, s)
	}
	dollars, err := balance(tokens).Dollars()
	if err != nil {
		return nil, err
	}
	cents, err := balance(tokens).Cents()
	if err != nil {
		return nil, err
	}
	n := dollars*100 + cents
	c := pb.Cent_builder{Value: &n}.Build()
	return c, nil
}

// balance is a pair of dollars and cents of ####.## value.
type balance []string

const (
	idxBalanceDollars = iota
	idxBalanceCents
)

// Dollars returns the integral part of the balance, the number before the
// dot, e.g. 123 in 123.45.
// It returns an error if the value is not a number.
func (b balance) Dollars() (int32, error) {
	n, err := strconv.Atoi(b[idxBalanceDollars])
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrDollar, err)
	}
	return int32(n), nil
}

// Cents parses the fractional part of the balance, the two digits after the
// dot, e.g. 45 in 123.45.
// It returns an error if cents are not a number of does not take two digits.
func (b balance) Cents() (int32, error) {
	s := b[idxBalanceCents]
	if len(s) != 2 {
		return 0, fmt.Errorf("%w: %s must be #.##", ErrCent, s)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrCent, err)
	}
	return int32(n), nil
}
