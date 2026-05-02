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
		Pos:  lex.Position{Start: 1, End: 5},
	}
	fmt.Println(t)
	// Output:
	// 1..5: number 12.3
}
