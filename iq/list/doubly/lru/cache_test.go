// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lru_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/list/doubly/lru"
)

func newCache(t *testing.T, capacity int, items ...lru.Item) *lru.Cache {
	t.Helper()
	cache, err := lru.NewCache(capacity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, item := range items {
		cache.Put(item.K, item.V)
	}
	return cache
}

func TestCachePut(t *testing.T) {
	tests := []struct {
		name  string
		cache *lru.Cache
		items []lru.Item
		want  []lru.Item
	}{
		{
			name:  "no evict on add one to empty cap two",
			cache: newCache(t, 2),
			items: []lru.Item{
				{1, 10},
			},
			want: []lru.Item{
				{1, 10},
			},
		},
		{
			name:  "no evict on add one to size one cap two",
			cache: newCache(t, 2, lru.Item{1, 10}),
			items: []lru.Item{
				{2, 20},
			},
			want: []lru.Item{
				{1, 10},
				{2, 20},
			},
		},
		{
			name:  "evict on add one to size two cap two",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			items: []lru.Item{
				{3, 30},
			},
			want: []lru.Item{
				{2, 20},
				{3, 30},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, item := range tc.items {
				tc.cache.Put(item.K, item.V)
			}

			if diff := cmp.Diff(tc.want, tc.cache.Items()); diff != "" {
				t.Errorf("lru.Cache state mistmatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	tests := []struct {
		name    string
		cache   *lru.Cache
		key     int
		wantVal int
		wantOk  bool
	}{
		{
			name:  "cache miss on empty",
			cache: newCache(t, 2),
		},
		// size 1 cap 2
		{
			name:    "cache hit on one item",
			cache:   newCache(t, 2, lru.Item{1, 10}),
			key:     1,
			wantVal: 10,
			wantOk:  true,
		},
		{
			name:  "cache miss on one item",
			cache: newCache(t, 2, lru.Item{1, 10}),
			key:   2,
		},
		// size 2 cap 2
		{
			name:    "cache hit on two items get first",
			cache:   newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:     1,
			wantVal: 10,
			wantOk:  true,
		},
		{
			name:    "cache hit on two items get second",
			cache:   newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:     2,
			wantVal: 20,
			wantOk:  true,
		},
		{
			name:  "cache miss on two items",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:   3,
		},
		{
			name:  "cache miss on two items get evicted",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}, lru.Item{3, 30}),
			key:   1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotVal, gotOk := tc.cache.Get(tc.key)

			if gotOk != tc.wantOk {
				t.Errorf("(*lru.Cache).Get(%v) = _, %v; want %v", tc.key, gotOk, tc.wantOk)
			}
			if gotVal != tc.wantVal {
				t.Errorf("(*lru.Cache).Get(%v) = %v, _; want %v", tc.key, gotVal, tc.wantVal)
			}
		})
	}
}

func TestCacheGetRebalance(t *testing.T) {
	tests := []struct {
		name  string
		cache *lru.Cache
		key   int
		want  []lru.Item
	}{
		{
			name:  "no rebalance on empty",
			cache: newCache(t, 2),
			key:   1,
		},
		// size 1 cap 2
		{
			name:  "no rebalance on one item",
			cache: newCache(t, 2, lru.Item{1, 10}),
			key:   1,
			want: []lru.Item{
				{1, 10},
			},
		},
		// size 2 cap 2
		{
			name:  "no rebalance on two items get recent",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:   2,
			want: []lru.Item{
				{1, 10},
				{2, 20},
			},
		},
		{
			name:  "rebalance on two items get least recent",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:   1,
			want: []lru.Item{
				{2, 20},
				{1, 10},
			},
		},
		{
			name:  "no rebalance on two items cache miss",
			cache: newCache(t, 2, lru.Item{1, 10}, lru.Item{2, 20}),
			key:   3,
			want: []lru.Item{
				{1, 10},
				{2, 20},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Get(tc.key)

			if diff := cmp.Diff(tc.want, tc.cache.Items()); diff != "" {
				t.Errorf("lru.Cache state mistmatch (-want, +got):\n%s", diff)
			}
		})
	}
}
