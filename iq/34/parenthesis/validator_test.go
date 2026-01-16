// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parenthesis_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/skhal/lab/iq/34/parenthesis"
)

func ExampleValidate() {
	fmt.Println(parenthesis.Validate("[()]"))
	// Output:
	// true
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		want   bool
	}{
		{
			name: "valid",
			inputs: []string{
				"",
				"()",
				"[]",
				"{}",
				"()()",
				"()[]",
				"(){}",
				"[]()",
				"[][]",
				"[]{}",
				"{}()",
				"{}[]",
				"{}{}",
				"(())",
				"([])",
				"({})",
				"[()]",
				"[[]]",
				"[{}]",
				"{()}",
				"{[]}",
				"{{}}",
			},
			want: true,
		},
		{
			name: "invalid",
			inputs: []string{
				"(",
				")",
				"[",
				"]",
				"{",
				"}",
				"(()",
				")()",
				"[()",
				"]()",
				"{()",
				"}()",
				"()(",
				"())",
				"()[",
				"()]",
				"(){",
				"()}",
				"([]",
				")[]",
				"[[]",
				"][]",
				"{[]",
				"}[]",
				"[](",
				"[])",
				"[][",
				"[]]",
				"[]{",
				"[]}",
				"({}",
				"){}",
				"[{}",
				"]{}",
				"{{}",
				"}{}",
				"{}(",
				"{})",
				"{}[",
				"{}]",
				"{}{",
				"{}}",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { testInputs(t, tc.inputs, tc.want) })
	}
}

func testInputs(t *testing.T, inputs []string, want bool) {
	t.Helper()
	for i, s := range inputs {
		s := s
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := parenthesis.Validate(s)

			if want != got {
				t.Errorf("parenthesis.Validate(%q) = %v; want %v", s, got, want)
			}
		})
	}
}
