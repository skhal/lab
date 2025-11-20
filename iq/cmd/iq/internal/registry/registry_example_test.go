// Copyright 2025 Samvel Khalatyan. All rights reserved.

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
