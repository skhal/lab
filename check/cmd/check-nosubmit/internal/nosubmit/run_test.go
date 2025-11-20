// Copyright 2025 Samvel Khalatyan. All rights reserved.

package nosubmit_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

func ExampleRun() {
	readFileFn := func(f string) ([]byte, error) {
		if f != "foo.txt" {
			return nil, fmt.Errorf("error opening file %s", f)
		}
		s := `
// DO NOT SUBMIT: work in progress
		`
		return []byte(s), nil
	}
	ctx := context.Background()
	cfg := &nosubmit.Config{
		ReadFileFn: readFileFn,
	}
	if err := nosubmit.Run(ctx, cfg, "foo.txt"); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// check error
}

func TestHasNoSubmit(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{name: "empty"},
		{
			name: "pass",
			data: `
test data
`,
		},
		{
			name: "nosubmit",
			data: `
test data
// DO NOT SUBMIT
`,
			want: true,
		},
		{
			name: "nosubmit with comment",
			data: `
test data
// DO NOT SUBMIT: description
`,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := nosubmit.Check([]byte(tc.data))

			if tc.want != got {
				t.Errorf("nosubmit.Check(...) = %v; want %v", got, tc.want)
				t.Logf("data:\n%s", tc.data)
			}
		})
	}
}
