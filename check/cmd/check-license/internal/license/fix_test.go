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
	wantC := `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`
	wantHTML := `<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->
`
	wantShell := `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
`
	tests := []struct {
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		// keep-sorted start
		{filename: ".bazelrc", holder: "Tester", want: wantShell},
		{filename: ".clangd", holder: "Tester", want: wantShell},
		{filename: "test", holder: "Tester", want: wantShell},
		{filename: "test.cc", holder: "Tester", want: wantC},
		{filename: "test.conf", holder: "Tester", want: wantShell},
		{filename: "test.ctags", holder: "Tester", want: wantShell},
		{filename: "test.go", holder: "Tester", want: wantC},
		{filename: "test.h", holder: "Tester", want: wantC},
		{filename: "test.html", holder: "Tester", want: wantHTML},
		{filename: "test.sh", holder: "Tester", want: wantShell},
		{filename: "test.txt", holder: "Tester", want: wantShell},
		{filename: "test.yaml", holder: "Tester", want: wantShell},
		// keep-sorted end
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
	wantC := `// Copyright 2026 Tester. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

test
`
	wantHTML := `<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

test
`
	wantShell := `# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

test
`
	tests := []struct {
		data     string
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		// keep-sorted start
		{data: "test\n", filename: ".bazelrc", holder: "Tester", want: wantShell},
		{data: "test\n", filename: ".clangd", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test.cc", holder: "Tester", want: wantC},
		{data: "test\n", filename: "test.conf", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test.ctags", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test.go", holder: "Tester", want: wantC},
		{data: "test\n", filename: "test.h", holder: "Tester", want: wantC},
		{data: "test\n", filename: "test.html", holder: "Tester", want: wantHTML},
		{data: "test\n", filename: "test.sh", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test.txt", holder: "Tester", want: wantShell},
		{data: "test\n", filename: "test.yaml", holder: "Tester", want: wantShell},
		// keep-sorted end
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
		name     string
		data     string
		filename string
		holder   string
		want     string
		wantErr  error
	}{
		{
			name: "no ext shebang and text",
			data: `#!/bin/sh
test
`,
			filename: "test",
			holder:   "Tester",
			want: `#!/bin/sh
#
# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

test
`,
		},
		{
			name: "no ext shebang only",
			data: `#!/bin/sh
`,
			filename: "test",
			holder:   "Tester",
			want: `#!/bin/sh
#
# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
`,
		},
		{
			name: "html doctype and text",
			data: `<!doctype html>
test
`,
			filename: "test.html",
			holder:   "Tester",
			want: `<!doctype html>
<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

test
`,
		},
		{
			name: "html doctype only",
			data: `<!doctype html>
`,
			filename: "test.html",
			holder:   "Tester",
			want: `<!doctype html>
<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->
`,
		},
		{
			name: "html doctype mixed case",
			data: `<!DoCTypE HtMl>
test
`,
			filename: "test_ignore_case.html",
			holder:   "Tester",
			want: `<!DoCTypE HtMl>
<!--
  Copyright 2026 Tester. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

test
`,
		},
		{
			name: "sh ext shebang and text",
			data: `#!/bin/sh
test
`,
			filename: "test.sh",
			holder:   "Tester",
			want: `#!/bin/sh
#
# Copyright 2026 Tester. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

test
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
