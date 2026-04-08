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
			name: "zero value",
			want: []byte{0x00, 0x00},
		},
		{
			name: "free",
			h:    heap.Header{Size: 5},
			want: []byte{0x00, 0x05},
		},
		{
			name: "free with allocated prev",
			h:    heap.Header{AllocatedPrev: true, Size: 5},
			want: []byte{1 << 6, 0x05},
		},
		{
			name: "allocated",
			h:    heap.Header{Allocated: true, Size: 5},
			want: []byte{1 << 7, 0x05},
		},
		{
			name: "allocated with allocated prev",
			h:    heap.Header{Allocated: true, AllocatedPrev: true, Size: 5},
			want: []byte{1<<7 | 1<<6, 0x05},
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
			name: "zero value",
			data: []byte{0x00, 0x00},
		},
		{
			name: "free",
			data: []byte{0x00, 0x05},
			want: heap.Header{Size: 5},
		},
		{
			name: "free with allocated prev",
			data: []byte{1 << 6, 0x05},
			want: heap.Header{AllocatedPrev: true, Size: 5},
		},
		{
			name: "allocated",
			data: []byte{1 << 7, 0x05},
			want: heap.Header{Allocated: true, Size: 5},
		},
		{
			name: "allocated with allocated prev",
			data: []byte{1<<7 | 1<<6, 0x05},
			want: heap.Header{Allocated: true, AllocatedPrev: true, Size: 5},
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

func TestFooter_Marshal(t *testing.T) {
	tests := []struct {
		name string
		f    heap.Footer
		want []byte
	}{
		{
			name: "zero value",
			want: []byte{0x00, 0x00},
		},
		{
			name: "size 5",
			f:    heap.Footer{Size: 5},
			want: []byte{0x00, 0x05},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.f.Marshal()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Marshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFooter_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want heap.Footer
	}{
		{
			name: "zero value",
			data: []byte{0x00, 0x00},
		},
		{
			name: "size 5",
			data: []byte{0x00, 0x05},
			want: heap.Footer{Size: 5},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := heap.Footer{}

			f.Unmarshal(tc.data)

			if diff := cmp.Diff(tc.want, f); diff != "" {
				t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
