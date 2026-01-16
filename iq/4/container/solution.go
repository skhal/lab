// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package container

import "math"

type Volume int

const VolumeEmpty Volume = 0

func Find(nn []int) Volume {
	if len(nn) < 2 {
		return VolumeEmpty
	}
	var (
		maxVolume Volume
		i1        = 0
		i2        = len(nn) - 1
	)
	for i1 < i2 {
		if volume := getVolume(nn, i1, i2); volume > maxVolume {
			maxVolume = volume
		}
		switch {
		case nn[i1] < nn[i2]:
			i1 += 1
		case nn[i1] > nn[i2]:
			i2 -= 1
		default:
			i1 += 1
			i2 -= 1
		}
	}
	return maxVolume
}

func getVolume(nn []int, i, j int) Volume {
	return Volume(min(nn[i], nn[j]) * (j - i))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
