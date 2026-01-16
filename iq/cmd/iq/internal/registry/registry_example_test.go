// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package registry_test

import (
	"fmt"

	"github.com/skhal/lab/iq/cmd/iq/internal/registry"
)

func ExampleLoad() {
	cfg := &registry.Config{File: "questions.txtpb"}
	_, err := registry.Load(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Output:
}
