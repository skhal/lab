// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stripzero

func CleanOptimal(m [][]int) {
	if m == nil {
		return
	}
	var (
		zero_first_row = false
		zero_first_col = false
	)
	for r, row := range m {
		for c, n := range row {
			if n != 0 {
				continue
			}
			if r == 0 {
				zero_first_row = true
			}
			if c == 0 {
				zero_first_col = true
			}
			m[r][0] = 0
			m[0][c] = 0
		}
	}
	for r, rlen := 1, len(m); r < rlen; r++ {
		if n := m[r][0]; n != 0 {
			continue
		}
		cleanRow(m, r)
	}
	for c, clen := 1, len(m[0]); c < clen; c++ {
		if n := m[0][c]; n != 0 {
			continue
		}
		cleanColumn(m, c)
	}
	if zero_first_row {
		cleanRow(m, 0)
	}
	if zero_first_col {
		cleanColumn(m, 0)
	}
}
