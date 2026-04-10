// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cycle demonstrates cycle detection in a directed graph.
package cycle

import (
	"errors"
	"strings"

	"github.com/skhal/lab/go/slices"
)

// ErrCycle means the graph has a cycle.
var ErrCycle = errors.New("cycle error")

// Node is a graph node with a human readable identifier and a list of edges
// to other nodes.
type Node struct {
	ID   string  // human readable identifier
	Deps []*Node // edges to other nodes
}

// HasCycle uses depth-first search to detect cycles in the graph. It reports
// first found cycle.
//
// The initial idea was brought from x/tools/go/analysis/validate.go[^1], where
// Validate() uses colors to:
//
//   - keep track of the nodes while parsing through the graph.
//
//   - build a cycle path by traversing nodes that are marked grey
//     (in-progress).
//
//   - ensure there are no duplicate root nodes. A sub-graph without cycles has
//     all nodes marked black. A unique root node is marked finished.
//     A duplicate root node would have finished-mark set twice.
//
// This code takes a different approach: it uses a single seen-map to keep
// track of visited nodes and a path-stack to report cycle, which hopefully is
// more readable for the following reasons:
//
//   - no cognitive load to keep track of color/node-state association (white,
//     grey, black, finished are not intuitive to associate with not-seen,
//     seen, no-cycles, unique-root node states correspondingly in Validate()).
//
//   - the cycle path is available right away for quick error reporting if the
//     node was visited in DFS, making error generation quick and easy.
//     Of course the logic can be encapsulated in a function, but current
//     version avoids it for readability.
//
// [^1]: https://cs.opensource.google/go/x/tools/+/master:go/analysis/validate.go;drc=d0d86e40a80dcab58f5cd2fa5f81e650d0777817 // NOLINT
func HasCycle(roots []*Node) error {
	var (
		seen = make(map[*Node]bool)
		path []*Node
	)
	var dfs func(*Node) error
	dfs = func(node *Node) error {
		path = append(path, node)
		if seen[node] {
			return CycleError(path)
		}
		seen[node] = true
		for _, next := range node.Deps {
			if err := dfs(next); err != nil {
				return err
			}
		}
		path = path[:len(path)-1]
		return nil
	}
	for _, node := range roots {
		if err := dfs(node); err != nil {
			return err
		}
	}
	return nil
}

// CycleError holds a cycle in the graph.
type CycleError []*Node

// Error implements [builin.error] interface.
func (err CycleError) Error() string {
	ids := slices.MapFunc(err, func(n *Node) string { return n.ID })
	return strings.Join(ids, ",")
}

// Is makes CycleError equivalent to [ErrCycle].
func (err CycleError) Is(other error) bool {
	return other == ErrCycle
}
