// Copyright 2025 Samvel Khalatyan. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s url\n", os.Args[0])
		os.Exit(1)
	}

	rawUrl := flag.Args()[0]
	escapedUrl := url.PathEscape(rawUrl)
	fmt.Println(escapedUrl)
}
