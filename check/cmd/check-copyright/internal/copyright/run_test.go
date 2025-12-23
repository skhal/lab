// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package copyright_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/check/cmd/check-copyright/internal/copyright"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		data string
		want error
	}{
		{
			name: "empty",
			want: copyright.ErrNotFound,
		},
		{
			name: "valid c style",
			data: `
// Copyright 2025 John Doe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
		},
		{
			name: "valid sh style",
			data: `
# Copyright 2025 John Doe. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
			`,
		},
		{
			name: "valid vim style",
			data: `
" Copyright 2025 John Doe. All rights reserved.
"
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
			`,
		},
		{
			name: "missing year",
			data: `
// Copyright John Doe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
			want: copyright.ErrNotFound,
		},
		{
			name: "missing author",
			data: `
// Copyright 2025. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
			want: copyright.ErrNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := copyright.Config{
				ReadFile: newReadFile(t, []byte(tc.data)),
			}

			err := copyright.Run(&cfg, "/nonexistent")

			if !errors.Is(err, tc.want) {
				t.Errorf("copyright.Run() unexpected error: %v; want %v", err, tc.want)
				t.Logf("data:\n%s", tc.data)
			}
		})
	}
}

func newReadFile(t *testing.T, data []byte) copyright.ReadFileFn {
	t.Helper()
	return func(string) ([]byte, error) {
		t.Helper()
		return data, nil
	}
}
