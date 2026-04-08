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

func TestEncoder_Encode(t *testing.T) {
	tests := []struct {
		name string
		size int
		h    heap.Header
		addr int
		want []byte
	}{
		{
			name: "zero value header",
			size: 5,
			addr: 2,
			want: []byte{0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "allocated block size 3",
			size: 5,
			h:    heap.Header{Allocated: true, Size: 3},
			addr: 2,
			want: []byte{1 << 7, 0x03, 0x00, 0x00, 0x00},
		},
		{
			name: "allocated block size 3 with allocated prev",
			size: 5,
			h:    heap.Header{Allocated: true, AllocatedPrev: true, Size: 3},
			addr: 2,
			want: []byte{1<<7 | 1<<6, 0x03, 0x00, 0x00, 0x00},
		},
		{
			name: "free block size 3",
			size: 5,
			h:    heap.Header{Size: 3},
			addr: 2,
			want: []byte{0x00, 0x03, 0x00, 0x00, 0x03},
		},
		{
			name: "free block size 3 with allocated prev",
			size: 5,
			h:    heap.Header{AllocatedPrev: true, Size: 3},
			addr: 2,
			want: []byte{1 << 6, 0x03, 0x00, 0x00, 0x03},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := make([]byte, tc.size)
			enc := heap.Encoder(buf)

			enc.Encode(&tc.h, tc.addr)

			if diff := cmp.Diff(tc.want, buf); diff != "" {
				t.Errorf("Encode(%v, %d) mismatch (-want +got):\n%s", tc.h, tc.addr, diff)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		addr int
		want heap.Header
	}{
		{
			name: "zero value header",
			data: []byte{0x00, 0x00, 0x00, 0x00, 0x00},
			addr: 2,
		},
		{
			name: "allocated block size 3",
			data: []byte{1 << 7, 0x03, 0x00, 0x00, 0x00},
			addr: 2,
			want: heap.Header{Allocated: true, Size: 3},
		},
		{
			name: "allocated block size 3 with allocated prev",
			data: []byte{1<<7 | 1<<6, 0x03, 0x00, 0x00, 0x00},
			addr: 2,
			want: heap.Header{Allocated: true, AllocatedPrev: true, Size: 3},
		},
		{
			name: "free block size 3",
			data: []byte{0x00, 0x03, 0x00, 0x00, 0x03},
			addr: 2,
			want: heap.Header{Size: 3},
		},
		{
			name: "free block size 3 with allocated prev",
			data: []byte{1 << 6, 0x03, 0x00, 0x00, 0x03},
			addr: 2,
			want: heap.Header{AllocatedPrev: true, Size: 3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dec := heap.Decoder(tc.data)
			var got heap.Header

			dec.Decode(&got, tc.addr)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Decode(_, %d) mismatch (-want +got):\n%s", tc.addr, diff)
			}
		})
	}
}

func TestDecoder_DecodePrevFooter(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		addr int
		want heap.Footer
	}{
		{
			name: "zero value footer",
			data: []byte{0x00, 0x00, 0x00, 0x01, 0x00},
			addr: 4,
		},
		{
			name: "block size 3",
			data: []byte{0x00, 0x03, 0x00, 0x01, 0x00},
			addr: 4,
			want: heap.Footer{Size: 3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dec := heap.Decoder(tc.data)
			var got heap.Footer

			dec.DecodePrevFooter(&got, tc.addr)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("DecodePrevFooter(_, %d) mismatch (-want +got):\n%s", tc.addr, diff)
			}
		})
	}
}
