// Copyright 2025 Samvel Khalatyan. All rights reserved.

package stripzero

func Clean(m [][]int) {
	if m == nil {
		return
	}
	cleanRows := make(map[int]struct{})
	cleanCols := make(map[int]struct{})
	for r, row := range m {
		for c, n := range row {
			if n != 0 {
				continue
			}
			cleanRows[r] = struct{}{}
			cleanCols[c] = struct{}{}
		}
	}
	for r := range cleanRows {
		cleanRow(m, r)
	}
	for c := range cleanCols {
		cleanColumn(m, c)
	}
}

func cleanRow(m [][]int, r int) {
	row := m[r]
	m[r] = make([]int, len(row))
}

func cleanColumn(m [][]int, c int) {
	for r := 0; r < len(m); r++ {
		m[r][c] = 0
	}
}
