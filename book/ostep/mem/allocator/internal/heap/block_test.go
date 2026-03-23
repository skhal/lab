// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/ostep/mem/allocator/internal/heap"
)

func TestHeader_Marshal(t *testing.T) {
	tests := []struct {
		name string
		h    heap.Header
		want []byte
	}{
		{
			name: "empty",
			want: []byte{0x00, 0x00},
		},
		{
			name: "free not empty",
			h:    heap.Header{Size: 5},
			want: []byte{0x00, 0x05},
		},
		{
			name: "allocated not empty",
			h:    heap.Header{Allocated: true, Size: 5},
			want: []byte{1 << 7, 0x05},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.h.Marshal()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Marshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want heap.Header
	}{
		{
			name: "zero data",
			data: []byte{0x00, 0x00},
		},
		{
			name: "free not empty",
			data: []byte{0x00, 0x05},
			want: heap.Header{Size: 5},
		},
		{
			name: "allocated not empty",
			data: []byte{1 << 7, 0x05},
			want: heap.Header{Allocated: true, Size: 5},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := heap.Header{}

			h.Unmarshal(tc.data)

			if diff := cmp.Diff(tc.want, h); diff != "" {
				t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
