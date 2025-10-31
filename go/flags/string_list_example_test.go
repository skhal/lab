// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flags_test

import (
	"flag"
	"fmt"

	"github.com/skhal/lab/go/flags"
)

func ExampleStringList() {
	var tags flags.StringList
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.Var(&tags, "tag", "comman separated tags")
	err := fs.Parse([]string{"-tag", "1", "-tag", "2,3", "-tag", ",,4"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tags)
	// Output:
	// [1 2 3 4]
}
