// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

// Address is a memory address.
type Address int

const (
	_  Address = 1 << (10 * iota)
	KB         // kilobyte
	MB         // megabyte
	GB         // gigabyte
)
