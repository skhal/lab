// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Check-nosubmit checks for presense of "DO NOT SUBMIT" comment.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/check/cmd/check-nosubmit/internal/nosubmit"
)

func main() {
	ctx := context.Background()
	cfg := &nosubmit.Config{
		ReadFileFn: os.ReadFile,
	}
	if err := nosubmit.Run(ctx, cfg, flag.Args()...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
