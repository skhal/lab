// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru_test

import (
	"fmt"

	"github.com/skhal/lab/iq/15/lru"
)

func ExampleCache_Put_evictsLeastRecentlyPutItem() {
	cache, _ := lru.NewCache(3)
	cache.Put(1, 100)
	cache.Put(2, 200)
	cache.Put(3, 300)
	cache.Put(4, 400)
	fmt.Println(cache)
	// Output:
	// [2:200 3:300 4:400]
}

func ExampleCache_Put_evictsLeastRecentlyUsedItem() {
	cache, _ := lru.NewCache(3)
	cache.Put(1, 100)
	cache.Put(2, 200)
	cache.Put(3, 300)
	cache.Get(1)
	cache.Put(4, 400)
	fmt.Println(cache)
	// Output:
	// [3:300 1:100 4:400]
}

func ExampleCache_Get_makesUsedItemRecent() {
	cache, _ := lru.NewCache(3)
	cache.Put(1, 100)
	cache.Put(2, 200)
	cache.Put(3, 300)
	cache.Get(1)
	fmt.Println(cache)
	// Output:
	// [2:200 3:300 1:100]
}
