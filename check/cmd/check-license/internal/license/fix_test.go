// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/check/cmd/check-license/internal/license"
)

func TestAdd_empty(t *testing.T) {
	tests := []struct {
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		{
			filename: "test",
			holder:   "Tester",
			want: `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
`,
		},
		{
			filename: "test.cc",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`,
		},
		{
			filename: "test.go",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`,
		},
		{
			filename: "test.h",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`,
		},
		{
			filename: "test.html",
			holder:   "Tester",
			want: `<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->
`,
		},
		{
			filename: "test.sh",
			holder:   "Tester",
			want: `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			got, err := license.Add([]byte(nil), tc.filename, tc.holder)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Add() got unexpected error %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("Add() mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestAdd_non_empty(t *testing.T) {
	tests := []struct {
		data     string
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		{
			data:     "echo test\n",
			filename: "test",
			holder:   "Tester",
			want: `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo test
`,
		},
		{
			data:     "namespace {}\n",
			filename: "test.cc",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

namespace {}
`,
		},
		{
			data:     "const test = 123\n",
			filename: "test.go",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

const test = 123
`,
		},
		{
			data:     "namespace {}\n",
			filename: "test.h",
			holder:   "Tester",
			want: `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

namespace {}
`,
		},
		{
			data:     "<br/>\n",
			filename: "test.html",
			holder:   "Tester",
			want: `<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

<br/>
`,
		},
		{
			data:     "echo test\n",
			filename: "test.sh",
			holder:   "Tester",
			want: `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo test
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			got, err := license.Add([]byte(tc.data), tc.filename, tc.holder)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Add() got unexpected error %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("Add() mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestAdd_header(t *testing.T) {
	tests := []struct {
		data     string
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		{
			data: `#!/bin/sh
echo test
`,
			filename: "test",
			holder:   "Tester",
			want: `#!/bin/sh
#
# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo test
`,
		},
		{
			data: `<!doctype html>
<br/>
`,
			filename: "test.html",
			holder:   "Tester",
			want: `<!doctype html>
<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

<br/>
`,
		},
		{
			data: `<!DOCTYPE HTML>
<br/>
`,
			filename: "test_ignore_case.html",
			holder:   "Tester",
			want: `<!DOCTYPE HTML>
<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

<br/>
`,
		},
		{
			data: `#!/bin/sh
echo test
`,
			filename: "test.sh",
			holder:   "Tester",
			want: `#!/bin/sh
#
# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

echo test
`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			got, err := license.Add([]byte(tc.data), tc.filename, tc.holder)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Add() got unexpected error %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("Add() mismatch (-want,+got):\n%s", diff)
			}
		})
	}
}
