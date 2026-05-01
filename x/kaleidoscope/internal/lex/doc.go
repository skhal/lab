// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lex implements lexer.
//
// EXAMPLE
//
//	# comment
//	def fib(x)
//		if x < 3 then
//			1
//		else
//			fib(x-1)+fib(x-2)
//
//	fib(10)
package lex
