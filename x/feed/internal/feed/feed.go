// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import "time"

type Feed <-chan *Item
type Item struct {
	Title     string
	Updated   *time.Time
	Published *time.Time
}
