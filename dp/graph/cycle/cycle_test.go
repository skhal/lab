// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cycle_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/skhal/lab/dp/graph/cycle"
	goslices "github.com/skhal/lab/go/slices"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		roots   []*cycle.Node
		want    error
		wantErr bool
	}{
		{
			name: "empty graph",
		},
		{
			name: "one node",
			roots: newGraph(t, map[string][]string{
				"A*": nil,
			}),
		},
		{
			name: "one node self cycle",
			roots: newGraph(t, map[string][]string{
				"A*": {"A"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "two connected nodes",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
			}),
		},
		{
			name: "two connected nodes with cycle",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"A"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "two connected nodes with self cycle A",
			roots: newGraph(t, map[string][]string{
				"A*": {"A", "B"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "two connected nodes with self cycle B",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"B"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "two disconnected nodes",
			roots: newGraph(t, map[string][]string{
				"A*": nil,
				"B*": nil,
			}),
		},
		{
			name: "two disconnected nodes with self cycle A",
			roots: newGraph(t, map[string][]string{
				"A*": {"A"},
				"B*": nil,
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "two disconnected nodes with self cycle B",
			roots: newGraph(t, map[string][]string{
				"A*": nil,
				"B*": {"B"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"C"},
			}),
		},
		{
			name: "three connected nodes with self cycle A",
			roots: newGraph(t, map[string][]string{
				"A*": {"A", "B"},
				"B":  {"C"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes with self cycle B",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"B", "C"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes with self cycle C",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"C"},
				"C":  {"C"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes with cycle C A",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"C"},
				"C":  {"A"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes with cycle B A",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"A", "C"},
			}),
			want: cycle.ErrCycle,
		},
		{
			name: "three connected nodes with cycle B C",
			roots: newGraph(t, map[string][]string{
				"A*": {"B"},
				"B":  {"C"},
				"C":  {"B"},
			}),
			want: cycle.ErrCycle,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cycle.HasCycle(tc.roots)

			if !errors.Is(err, tc.want) {
				t.Errorf("Validate() unsexpected error %v, want %v", err, tc.want)
				logGraph(t, tc.roots)
			}
		})
	}
}

func newGraph(t *testing.T, edges map[string][]string) []*cycle.Node {
	t.Helper()
	var (
		nodes = make(map[string]*cycle.Node)
		roots []*cycle.Node
	)
	node := func(id string) *cycle.Node {
		id, root := strings.CutSuffix(id, "*")
		n, ok := nodes[id]
		if !ok {
			n = &cycle.Node{ID: id}
			nodes[id] = n
		}
		// the edges is a map with random order: a root node can be created
		// through an edge.
		if root && !slices.Contains(roots, n) {
			roots = append(roots, n)
		}
		return n
	}
	for id, deps := range edges {
		n := node(id)
		for _, id := range deps {
			dep := node(id)
			n.Deps = append(n.Deps, dep)
		}
	}
	return roots
}

func logGraph(t *testing.T, roots []*cycle.Node) {
	t.Helper()
	var (
		seen  = make(map[*cycle.Node]bool)
		edges []string
	)
	var dfs func(format string, node *cycle.Node)
	dfs = func(format string, node *cycle.Node) {
		if seen[node] {
			return
		}
		seen[node] = true
		deps := goslices.MapFunc(node.Deps, func(n *cycle.Node) string {
			return n.ID
		})
		str := fmt.Sprintf(format, node.ID, strings.Join(deps, ","))
		edges = append(edges, str)
		for _, n := range node.Deps {
			dfs("%s: %s", n)
		}
	}
	for _, node := range roots {
		dfs("%s*: %s", node)
	}
	switch len(edges) {
	case 0:
		t.Log("graph: <empty>")
	default:
		t.Logf("graph:\n%s", strings.Join(edges, "\n"))
	}
}
