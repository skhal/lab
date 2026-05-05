// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex_test

import (
	"fmt"

	"github.com/skhal/lab/x/kscope/internal/lex"
)

func ExampleToken_String() {
	t := lex.Token{
		Kind: lex.TokNum,
		Val:  "12.3",
	}
	fmt.Println(t)
	// Output:
	// number "12.3"
}
