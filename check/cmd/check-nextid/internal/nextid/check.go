// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nextid

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/bufbuild/protocompile/ast"
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
)

// ErrNextID means that next is incorrect or failed to parse.
var ErrNextID = errors.New("invalid next id")

// ErrRange means that it failed to process a reserved range.
var ErrRange = errors.New("invalid range")

// CheckFile validates that a .proto file parses and has next-id comments set
// to the (last-id + 1).
func CheckFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	fn, err := parser.Parse(name, f, reporter.NewHandler(nil))
	if err != nil {
		return err
	}
	err = CheckFileNode(fn)
	if err != nil {
		return err
	}
	return nil
}

// CheckFileNode validates a .proto AST's top-level declarations.
func CheckFileNode(file *ast.FileNode) error {
	for _, d := range file.Decls {
		if err := checkFileElement(file, d); err != nil {
			return err
		}
	}
	return nil
}

func checkFileElement(file *ast.FileNode, elem ast.FileElement) error {
	switch n := elem.(type) {
	case *ast.EnumNode:
		if err := checkEnumNode(file, n); err != nil {
			return err
		}
	case *ast.MessageNode:
		if err := checkMessageNode(file, n); err != nil {
			return err
		}
	}
	return nil
}

type enumNode struct {
	node *ast.EnumNode
}

// Node returns wrapped node.
func (n *enumNode) Node() ast.CompositeNode {
	return n.node
}

// Name returns wrapped node's name.
func (n *enumNode) Name() string {
	return n.node.Name.Val
}

// Kind returns "enum" for kind.
func (n *enumNode) Kind() string {
	return "enum"
}

type messageNode struct {
	node *ast.MessageNode
}

// Node returns wrapped node.
func (n *messageNode) Node() ast.CompositeNode {
	return n.node
}

// Name returns wrapped node's name.
func (n *messageNode) Name() string {
	return n.node.Name.Val
}

// Kind returns "message" for kind.
func (n *messageNode) Kind() string {
	return "message"
}

type compositeNode interface {
	Node() ast.CompositeNode
	Name() string
	Kind() string
}

func checkEnumNode(file *ast.FileNode, enum *ast.EnumNode) error {
	n := enumNode{enum}
	return checkLeadingComments(&n, file.NodeInfo(enum).LeadingComments())
}

func checkMessageNode(file *ast.FileNode, msg *ast.MessageNode) error {
	n := messageNode{msg}
	err := checkLeadingComments(&n, file.NodeInfo(msg).LeadingComments())
	if err != nil {
		return err
	}
	for _, child := range msg.Children() {
		switch n := child.(type) {
		case *ast.EnumNode:
			err := checkEnumNode(file, n)
			if err != nil {
				return err
			}
		case *ast.MessageNode:
			err := checkMessageNode(file, n)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkLeadingComments(cn compositeNode, comms ast.Comments) error {
	// Scan for the next-id comment from the end, where it is likely to be found.
	for i := comms.Len() - 1; i >= 0; i-- {
		comm := comms.Index(i)
		switch ok, err := checkComment(cn, comm); {
		case err != nil:
			return fmt.Errorf("%s %s %s: %w", comm.Start(), cn.Kind(), cn.Name(), err)
		case !ok:
			continue
		}
		break
	}
	return nil
}

func checkComment(cn compositeNode, comm ast.Comment) (ok bool, err error) {
	nextid, ok, err := ParseNextID(comm.RawText())
	if err != nil {
		return true, err
	}
	if !ok {
		return false, nil
	}
	lastid, err := getLastID(cn.Node())
	if err != nil {
		return true, fmt.Errorf("failed to extract last id: %w", err)
	}
	if nextid != lastid+1 {
		return true, fmt.Errorf("last id %d: next id %d: %w", lastid, nextid, ErrNextID)
	}
	return true, nil
}

var reNextID = regexp.MustCompile(`(?i)^\s*// next id: (.+)$`)

// ParseNextID parses next-id value from the comment line using
// [strconv.ParseUint].
func ParseNextID(text string) (nextid uint64, ok bool, err error) {
	matches := reNextID.FindStringSubmatch(text)
	if matches == nil {
		return
	}
	nextid, err = strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("%w %q", ErrNextID, matches[1])
		return
	}
	if nextid == 0 {
		return 0, false, fmt.Errorf("%w: must be positive", ErrNextID)
	}
	return nextid, true, nil
}

func getLastID(node ast.CompositeNode) (uint64, error) {
	var lastid uint64
	for _, child := range node.Children() {
		switch n := child.(type) {
		case *ast.EnumValueNode:
			id, ok := n.Number.AsUint64()
			if !ok {
				return 0, fmt.Errorf("enumerator %s: invalid number %v", n.Name.Val, n.Number)
			}
			if id > lastid {
				lastid = id
			}
		case *ast.FieldNode:
			if id := n.Tag.Val; id > lastid {
				lastid = id
			}
		case *ast.ReservedNode:
			for _, r := range n.Ranges {
				id, err := getRangeLastID(r)
				if err != nil {
					return 0, err
				}
				if id > lastid {
					lastid = id
				}
			}
		}
	}
	return lastid, nil
}

func getRangeLastID(rn *ast.RangeNode) (uint64, error) {
	n, ok := rn.StartVal.AsUint64()
	if !ok {
		return 0, ErrRange
	}
	if rn.To == nil {
		return n, nil
	}
	if rn.Max != nil {
		// Assuming UINT32_MAX
		return math.MaxUint32, nil
	}
	n, ok = rn.EndVal.AsUint64()
	if !ok {
		return 0, ErrRange
	}
	return n, nil
}
