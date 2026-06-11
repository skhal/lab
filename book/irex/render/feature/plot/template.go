// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

var (
	//go:embed static
	efs embed.FS

	tmplPlotFeature = template.Must(func() (*template.Template, error) {
		// We want to indent nested template output to align HTML:
		//
		//		{{define "foo" -}}
		//		<div id="inner">
		//			demo
		//		</div>
		//		{{- end}}
		// 		<div id="container">
		//			{{template "foo"}}
		//		</div>
		//
		// The result is:
		//
		//		<div id="container">
		//			<div id="inner">
		//		  demo
		//    </div>
		//    </div>
		//
		// instead of:
		//
		//		<div id="container">
		//			<div id="inner">
		//		    demo
		//      </div>
		//    </div>
		//
		// Since {{template ...}} is an action, it can't be part of a pipeline.
		// The following code is illegal (assuming indent-function):
		//
		//		{{template "foo" | indent "\t"}}
		//
		// A solution is to introduce a function to replace the template action,
		// say include-function. Now it can be used in a pipeline to indent the
		// result of the nested template:
		//
		//		{{include "foo" | indent "\t"}}
		//
		// Ref: https://stackoverflow.com/questions/43821989/how-to-indent-content-of-included-template // NOLINT
		t := template.New("index.html")
		fmap := template.FuncMap{
			"include": func(name string, data any) (template.HTML, error) {
				var b bytes.Buffer
				if err := t.ExecuteTemplate(&b, name, data); err != nil {
					return "", err
				}
				return template.HTML(b.String()), nil
			},
			"indent": func(prefix string, s template.HTML) template.HTML {
				var (
					b         strings.Builder
					addPrefix bool
				)
				for l := range strings.Lines(string(s)) {
					if !addPrefix {
						fmt.Fprint(&b, l)
						addPrefix = true
					} else {
						fmt.Fprintf(&b, "%s%s", prefix, l)
					}
				}
				return template.HTML(b.String())
			},
		}
		return t.Funcs(fmap).ParseFS(efs, "static/index.html")
	}())

	tmplsPlotFeatureJS = template.Must(template.New("init.js").ParseFS(efs, "static/init.js"))
)
