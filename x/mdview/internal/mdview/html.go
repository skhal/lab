// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"context"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type (
	Markdown []byte // Markdown data
	HTML     []byte // HTML data
)

// ToHTML converts Markdown buffer to HTML.
func ToHTML(ctx context.Context, md Markdown) (HTML, error) {
	ast := parse(md)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return render(ast), nil
}

func parse(md Markdown) ast.Node {
	parser := parser.NewWithExtensions(parser.CommonExtensions)
	return parser.Parse(md)
}

func render(node ast.Node) HTML {
	opts := html.RendererOptions{
		Flags: html.CommonFlags,
	}
	renderer := html.NewRenderer(opts)
	return markdown.Render(node, renderer)
}
