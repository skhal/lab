// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package feed implements streaming interface to access RSS, Atom, etc. feeds.
package feed

import "time"

// Item is a feed item, be it RSS, Atom, etc.
type Item struct {
	Title     string     // Title of the item.
	Updated   *time.Time // When the item was updated.
	Published *time.Time // When the item was published.
}
