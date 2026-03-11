// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package issue verifies that commit message, represented by a file, has a
// reference to the issue, i.e., either `NO_ISSUE` or `Issue #123` line,
// ignoring case. The issue number can also have `user/repo#123` form.
package issue
