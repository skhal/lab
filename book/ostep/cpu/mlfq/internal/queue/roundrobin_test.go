// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/ostep/cpu/mlfq/internal/queue"
)

func TestRoundRobin_Append(t *testing.T) {
	tests := []struct {
		name    string
		items   []int
		wantLen int
	}{
		{
			name: "empty",
		},
		{
			name:    "one item",
			items:   []int{1},
			wantLen: 1,
		},
		{
			name:    "two items",
			items:   []int{1, 2},
			wantLen: 2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := new(queue.RoundRobin)
			for _, v := range tc.items {
				rr.Append(v)
			}

			if got := rr.Len(); got != tc.wantLen {
				t.Errorf("(*RoundRobin).Len() = %d, want %d", got, tc.wantLen)
			}
		})
	}
}

func TestRoundRobin_Append_with_next_pop(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		calls func(*queue.RoundRobin)
		want  []int
	}{
		{
			name:  "one item one next one append",
			items: []int{1},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Append(2)
			},
			want: []int{2, 1},
		},
		{
			name:  "one item one pop one append",
			items: []int{1},
			calls: func(rr *queue.RoundRobin) {
				rr.Pop()
				rr.Append(2)
			},
			want: []int{2},
		},
		{
			name:  "one item one next one pop one append",
			items: []int{1},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Pop()
				rr.Append(2)
			},
			want: []int{2},
		},
		{
			name:  "two items one next one append",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Append(3)
			},
			want: []int{2, 3, 1},
		},
		{
			name:  "two items two next one append",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
				rr.Append(3)
			},
			want: []int{3, 1, 2},
		},
		{
			name:  "two items one pop one append",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Pop()
				rr.Append(3)
			},
			want: []int{2, 3},
		},
		{
			name:  "two items two pop one append",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Pop()
				rr.Pop()
				rr.Append(3)
			},
			want: []int{3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)
			tc.calls(rr)

			var got []int
			for range tc.want {
				got = append(got, rr.Next().(int))
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("call (*RoundRobin).Next() %d times mismatch (-want +got):\n%s", len(tc.want), diff)
			}
		})
	}
}

func TestRoundRobin_Next(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  []int
	}{
		{
			name: "empty",
		},
		{
			name:  "one item",
			items: []int{1},
			want:  []int{1},
		},
		{
			name:  "one item one cycle",
			items: []int{1},
			want:  []int{1, 1},
		},
		{
			name:  "one item two cycles",
			items: []int{1},
			want:  []int{1, 1, 1},
		},
		{
			name:  "two items",
			items: []int{1, 2},
			want:  []int{1, 2},
		},
		{
			name:  "two items one cycle",
			items: []int{1, 2},
			want:  []int{1, 2, 1, 2},
		},
		{
			name:  "two items two cycles",
			items: []int{1, 2},
			want:  []int{1, 2, 1, 2, 1, 2},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)

			var got []int
			for range tc.want {
				got = append(got, rr.Next().(int))
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("call (*RoundRobin).Next() %d times mismatch (-want +got):\n%s", len(tc.want), diff)
			}
		})
	}
}

func TestRoundRobin_Next_panic(t *testing.T) {
	var err error
	rr := new(queue.RoundRobin)

	func() {
		defer func() {
			x := recover()
			if x == nil {
				return
			}
			if v, ok := x.(error); ok {
				err = v
			}
		}()
		rr.Next()
	}()

	if want := queue.ErrEmpty; !errors.Is(err, want) {
		t.Errorf("(*RoundRobin).Next() unexpected error %v; want %v", err, want)
	}
}

func TestRoundRobin_NextFunc(t *testing.T) {
	tests := []struct {
		name   string
		items  []int
		f      func(any) bool
		want   int
		wantOk bool
	}{
		{
			name:   "one item hit",
			items:  []int{1},
			f:      func(x any) bool { return x.(int) == 1 },
			want:   1,
			wantOk: true,
		},
		{
			name:  "one item miss",
			items: []int{1},
			f:     func(any) bool { return false },
		},
		{
			name:   "two items hit first",
			items:  []int{1, 2},
			f:      func(x any) bool { return x.(int) == 1 },
			want:   1,
			wantOk: true,
		},
		{
			name:   "two items hit second",
			items:  []int{1, 2},
			f:      func(x any) bool { return x.(int) == 2 },
			want:   2,
			wantOk: true,
		},
		{
			name:  "two items miss",
			items: []int{1, 2},
			f:     func(any) bool { return false },
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)

			got, ok := rr.NextFunc(tc.f)

			if ok != tc.wantOk {
				t.Errorf("NextFunc() = _, %v; want %v", ok, tc.wantOk)
			}
			if n := got.(int); n != tc.want {
				t.Errorf("NextFunc() = %d, _; want %d", got, tc.want)
			}
		})
	}
}

func TestRoundRobin_Pop(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		callNext int
		want     []int
	}{
		{
			name: "empty",
		},
		{
			name:  "one item",
			items: []int{1},
			want:  []int{1},
		},
		{
			name:     "one item one next",
			items:    []int{1},
			callNext: 1,
			want:     []int{1},
		},
		{
			name:  "two items",
			items: []int{1, 2},
			want:  []int{1, 2},
		},
		{
			// equivalent to never called Next()
			name:     "two items one next",
			items:    []int{1, 2},
			callNext: 1,
			want:     []int{1, 2},
		},
		{
			name:     "two items two next",
			items:    []int{1, 2},
			callNext: 2,
			want:     []int{2, 1},
		},
		{
			name:     "two items three next",
			items:    []int{1, 2},
			callNext: 3,
			want:     []int{1, 2},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)
			for range tc.callNext {
				rr.Next()
			}

			var got []int
			for range tc.want {
				got = append(got, rr.Pop().(int))
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("call (*RoundRobin).Pop() %d times mismatch (-want +got):\n%s", len(tc.want), diff)
			}
		})
	}
}

func TestRoundRobin_Pop_panic(t *testing.T) {
	tests := []struct {
		name  string
		items []int
	}{
		{
			name: "empty",
		},
		{
			name:  "one item",
			items: []int{1},
		},
		{
			name:  "two items",
			items: []int{1, 2},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)
			for range len(tc.items) {
				rr.Pop()
			}
			var err error

			func() {
				defer func() {
					x := recover()
					if x == nil {
						return
					}
					if v, ok := x.(error); ok {
						err = v
					}
				}()
				rr.Pop()
			}()

			if want := queue.ErrEmpty; !errors.Is(err, want) {
				t.Errorf("(*RoundRobin).Pop() unexpected error %v; want %v", err, want)
			}
		})
	}
}

func TestRoundRobin_Pop_with_next(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		calls func(*queue.RoundRobin)
		want  []int
	}{
		{
			name:  "one item one next",
			items: []int{1},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
			},
			want: []int{1},
		},
		{
			name:  "two items one next",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
			},
			want: []int{1, 2},
		},
		{
			name:  "two items two next",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
			},
			want: []int{2, 1},
		},
		{
			name:  "two items three next",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
				rr.Next()
			},
			want: []int{1, 2},
		},
		{
			name:  "three items one next",
			items: []int{1, 2, 3},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
			},
			want: []int{1, 2, 3},
		},
		{
			name:  "three items two next",
			items: []int{1, 2, 3},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
			},
			want: []int{2, 1, 3},
		},
		{
			name:  "three items three next",
			items: []int{1, 2, 3},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
				rr.Next()
			},
			want: []int{3, 2, 1},
		},
		{
			name:  "three items four next",
			items: []int{1, 2, 3},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
				rr.Next()
				rr.Next()
			},
			want: []int{1, 3, 2},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)
			tc.calls(rr)

			var got []int
			for range tc.want {
				got = append(got, rr.Pop().(int))
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("call (*RoundRobin).Pop() %d times mismatch (-want +got):\n%s", len(tc.want), diff)
			}
		})
	}
}

func TestRoundRobin_mix_next_pop(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		calls func(*queue.RoundRobin)
		want  []int
	}{
		{
			name:  "two items pop first",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Pop()
			},
			want: []int{2},
		},
		{
			name:  "two items pop second",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Next()
				rr.Next()
				rr.Pop()
			},
			want: []int{1},
		},
		{
			name:  "two items pop first",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Pop()
				rr.Next()
			},
			want: []int{2},
		},
		{
			name:  "two items pop all",
			items: []int{1, 2},
			calls: func(rr *queue.RoundRobin) {
				rr.Pop()
				rr.Next()
				rr.Pop()
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := newRoundRobin(t, tc.items)
			tc.calls(rr)

			var got []int
			for range tc.want {
				got = append(got, rr.Next().(int))
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("call (*RoundRobin).Next() %d times mismatch (-want +got):\n%s", len(tc.want), diff)
			}
		})
	}
}

func newRoundRobin(t *testing.T, nn []int) *queue.RoundRobin {
	t.Helper()
	rr := new(queue.RoundRobin)
	for _, v := range nn {
		rr.Append(v)
	}
	return rr
}
