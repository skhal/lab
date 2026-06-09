// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package render dispatches page rendering to feature renderers and assebles
// the final result into the finance results page.
package render

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/render/feature/plot"
	"github.com/skhal/lab/book/irex/web/queryparam"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrFeature means the feature has non-matching renderer registered.
	ErrFeature = errors.New("invalid feature")

	// ErrNoRenderer means the passed feature has no registered renderer.
	ErrNoRenderer = errors.New("no feature renderer")

	// ErrEmptyFeature means passed feature contains no feature extension.
	ErrEmptyFeature = errors.New("empty feature")
)

var (
	//go:embed static
	efs embed.FS

	tmplPage = template.Must(template.New("index.html").ParseFS(efs, "static/*.html"))
)

// Renderer is common interface of feature renderers.
type Renderer interface {
	// Render renders the feature into HTML. It returns an error if rendering
	// fails.
	Render() (template.HTML, error)
}

type renderFunc func(proto.Message) (Renderer, error)

var featureRenderers = map[protoreflect.ExtensionType]renderFunc{
	pb.E_PlotFeature_PlotFeature: dispatch(plot.NewRenderer),
}

// Render renders features included in the page. It preserves the features
// order.
func Render(p *pb.Page, w http.ResponseWriter, req *http.Request) error {
	type Header struct {
		Search string
	}
	d := struct {
		Header    Header
		Renderers []Renderer
	}{
		Renderers: make([]Renderer, len(p.GetFeatures())),
	}
	if q, ok := queryparam.Query(req); ok {
		d.Header.Search = q
	}
	for idx, feature := range p.GetFeatures() {
		fr, err := newFeatureRenderer(feature)
		if err != nil {
			return err
		}
		d.Renderers[idx] = fr
	}
	return tmplPage.Execute(w, d)
}

// newFeatureRenderer calls feature renderer, registered using feature
// extension.
// It returns an error if no renderer is registered for the feature and
// propagates feature renderer errors.
func newFeatureRenderer(feature *pb.Feature) (rend Renderer, err error) {
	f := func(ext protoreflect.ExtensionType, val any) bool {
		render, ok := featureRenderers[ext]
		if !ok {
			err = fmt.Errorf("%w: %s", ErrNoRenderer, ext)
			return false
		}
		v := val.(proto.Message)
		rend, err = render(v)
		return false // stop on the first feature
	}
	proto.RangeExtensions(feature, f)
	if err == nil && rend == nil {
		err = ErrEmptyFeature
	}
	return
}

type renderFeatureFunc[T proto.Message, R Renderer] func(msg T) R

// dispatch adapts feature-render function with strong type in function
// parameters to renderFunc, which takes proto.Message interface. It takes
// care of casting proto.Message to the feature type.
//
// The function returns an error if the cast fails, which should not happen
// as long as the mapping between the feature-render function and the extension
// match, i.e. the same message type and extension.
func dispatch[T proto.Message, R Renderer](f renderFeatureFunc[T, R]) renderFunc {
	return func(msg proto.Message) (Renderer, error) {
		feature, ok := msg.(T)
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrFeature, msg)
		}
		var r Renderer = f(feature)
		return r, nil
	}
}
