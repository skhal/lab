// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fin

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrFormat indicates an error in the string format when parsing a
	// Transaction from the reader.
	ErrFormat = errors.New("invalid format")
)

const dateFormat = "1/02/2006"

// Transaction represents a financial transaction.
type Transaction struct {
	// Customer is the customer name.
	Customer string

	// Date is when the transaction happens. Only use the date, no time.
	Date time.Time

	// Amount of the transaction up to hundreths precision.
	Amount float64
}

// String implements fmt.Stringer() interface.
func (t *Transaction) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s %s %.2f", t.Customer, t.Date.Format(dateFormat), t.Amount)
}

// Scan implements fmt.Scan() interface.
func (t *Transaction) Scan(state fmt.ScanState, verb rune) error {
	var customer string
	_, err := fmt.Fscan(state, &customer)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFormat, err)
	}
	tok, err := state.Token(true, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFormat, err)
	}
	date, err := time.Parse(dateFormat, string(tok))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFormat, err)
	}
	var amount float64
	_, err = fmt.Fscan(state, &amount)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFormat, err)
	}
	t.Customer = customer
	t.Date = date
	t.Amount = amount
	return nil
}
