// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sheet

import "fmt"

// Calculator drives cell calculatorion. It keeps track of visited cells to
// detect cycles and provides a mechanism to calculate cell references.
type Calculator struct {
	s    *Sheet
	seen map[string]bool
	path []string // references path for cycle detection
}

func newCalculator(s *Sheet) *Calculator {
	return &Calculator{s: s, seen: make(map[string]bool)}
}

// Calculate calculates cells. It returns an error if calculation fails for
// any of the cells.
func (calc *Calculator) Calculate() error {
	for id, c := range calc.s.data {
		if err := calc.calculate(id, c); err != nil {
			return fmt.Errorf("%w %s: calculate: %s", ErrCell, id, err)
		}
	}
	return nil
}

func (calc *Calculator) calculate(id string, c *cell) error {
	if c.calculated {
		return nil
	}
	calc.seen[id] = true
	calc.path = append(calc.path, id)
	defer func() {
		calc.path = calc.path[:len(calc.path)-1]
		calc.seen[id] = false
	}()
	res, err := calc.s.eng.Calculate(c.ir, calc.CalculateReference)
	if err != nil {
		return fmt.Errorf("%s: %s", id, err)
	}
	c.result = res
	c.calculated = true
	return nil
}

// CalculateReference calculates a reference value. It returns an error if
// calculation fails or reference calculator detects a circular dependency.
func (calc *Calculator) CalculateReference(id string) (float64, error) {
	if calc.seen[id] {
		// circular dependency
		return 0, fmt.Errorf("circular dependency - %v", calc.path)
	}
	c, ok := calc.s.data[id]
	if !ok {
		return 0, nil
	}
	if err := calc.calculate(id, c); err != nil {
		return 0, err
	}
	return c.result, nil
}
