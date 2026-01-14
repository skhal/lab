// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/check/cmd/check-license/internal/license"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		data string
		want error
	}{
		{
			name: "empty",
			want: license.ErrNotFound,
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
			name: "valid html style",
			data: `
<!--
  Copyright 2025 John Doe. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->`,
		},
		{
			name: "missing year",
			data: `
// Copyright John Doe. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
			want: license.ErrInvalid,
		},
		{
			name: "missing author",
			data: `
// Copyright 2025. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
			want: license.ErrInvalid,
		},
		{
			name: "mixed comment prefix",
			data: `
// Copyright 2025. All rights reserved.
#
" Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
			`,
			want: license.ErrInvalid,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := license.Check([]byte(tc.data))

			if !errors.Is(err, tc.want) {
				t.Errorf("license.Run() unexpected error: %v; want %v", err, tc.want)
				t.Logf("data:\n%s", tc.data)
			}
		})
	}
}
