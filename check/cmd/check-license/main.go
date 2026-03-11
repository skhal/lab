// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check-license verifies that the file includes a copyright statement.
package main

import (
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-license/internal/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
