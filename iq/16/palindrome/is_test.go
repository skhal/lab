// Copyright 2025 Samvel Khalatyan. All rights reserved.

package palindrome_test

import (
	"testing"

	"github.com/skhal/lab/iq/16/palindrome"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		list *palindrome.Node
		want bool
	}{
		{
			name: "nil",
			want: true,
		},
		{
			name: "empty",
			list: palindrome.NewList(),
			want: true,
		},
		{
			name: "one item",
			list: palindrome.NewList(1),
			want: true,
		},
		{
			name: "two items not palindrome",
			list: palindrome.NewList(1, 2),
		},
		{
			name: "two items palindrome",
			list: palindrome.NewList(1, 1),
			want: true,
		},
		{
			name: "three items not palindrome",
			list: palindrome.NewList(1, 2, 3),
		},
		{
			name: "three items palindrome",
			list: palindrome.NewList(1, 2, 1),
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := palindrome.Is(tc.list)

			if got != tc.want {
				t.Errorf("palindrome.Is(%s) = %v; want %v", tc.list, got, tc.want)
			}
		})
	}
}
