// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags_test

import (
	"flag"
	"fmt"

	"github.com/skhal/lab/x/go/flags"
)

func ExampleRequiredString() {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	foo := new(flags.RequiredString)
	bar := new(flags.RequiredString)
	fs.Var(foo, "foo", "foo is a required string")
	fs.Var(foo, "bar", "bar is a required string")
	fs.Parse([]string{"-foo", "foo-value"})
	fmt.Println("foo:", foo)
	fmt.Println("bar:", bar)
	// Output:
	// foo: foo-value
	// bar: <not parsed>
}
