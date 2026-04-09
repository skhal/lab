// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"fmt"
	"strconv"
	"strings"
)

type replayer struct {
	ops []string
	op  operation

	malloc OpFunc
	free   OpFunc
}

// OpFunc generates a malloc or free operations. The meaning of the integer
// parameter depends on the operation being generated: malloc(n) means generate
// a malloc operation to allocated n-bytes, free(n) means generate a free
// operation to release n-th allocation available at the moment of running the
// free operation.
type OpFunc func(int) operation

func newReplayer(ops []string, malloc, free OpFunc) *replayer {
	return &replayer{
		ops:    ops,
		malloc: malloc,
		free:   free,
	}
}

// Next moves the replayer to the next operation. It returns false if there are
// no operations left to replay.
func (rep *replayer) Next() bool {
	if len(rep.ops) == 0 {
		return false
	}
	rep.next()
	rep.ops = rep.ops[1:]
	return true
}

func (rep *replayer) next() {
	op := rep.ops[0]
	switch {
	case strings.HasPrefix(op, "+"):
		s := strings.TrimLeft(op, "+")
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		rep.op = rep.malloc(n)
	case strings.HasPrefix(op, "-"):
		s := strings.TrimLeft(op, "-")
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		rep.op = rep.free(n)
	default:
		panic(fmt.Errorf("invalid operation %q", op))
	}
}

// Op returns last operation produces by call to [replayer.Next].
func (rep *replayer) Op() operation {
	return rep.op
}
