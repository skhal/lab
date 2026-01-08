// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type (
	Markdown []byte // Markdown data
	HTML     []byte // HTML data
)

// ToHTML converts Markdown buffer to HTML.
func Render(md Markdown) HTML {
	parser := parser.NewWithExtensions(parser.CommonExtensions)
	ast := parser.Parse(md)
	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	return markdown.Render(ast, renderer)
}
