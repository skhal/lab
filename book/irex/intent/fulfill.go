// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package intent

import (
	"errors"
	"fmt"

	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrUnsupportedIntent means the intent does not have associated fulfiller.
	ErrUnsupportedIntent = errors.New("unsupported intent")

	// ErrIntent means invalid intent was passed to the fulfiller.
	ErrIntent = errors.New("invalid intent")
)

type fulfillFunc func(msg proto.Message) (proto.Message, error)

type fulfillInfo struct {
	f   fulfillFunc
	ext protoreflect.ExtensionType
}

var fulfillInfos = map[protoreflect.ExtensionType]fulfillInfo{
	pb.E_PlotIntent_PlotIntent: {
		f:   dispatch(fulfillPlotIntent),
		ext: pb.E_PlotFeature_PlotFeature,
	},
	pb.E_PingIntent_PingIntent: {
		f:   dispatch(fulfillPingIntent),
		ext: pb.E_PingFeature_PingFeature,
	},
}

// Fulfill dispatches the fulfillment to the intent fulfiller using the intent
// extension.
func Fulfill(intent *pb.Intent) (page *pb.Page, err error) {
	f := func(ext protoreflect.ExtensionType, val any) bool {
		info, ok := fulfillInfos[ext]
		if !ok {
			err = fmt.Errorf("%w: %s", ErrUnsupportedIntent, ext)
			return false
		}
		v := val.(proto.Message)

		var res proto.Message
		res, err = info.f(v)
		if err != nil {
			return false
		}

		feature := pb.Feature_builder{}.Build()
		proto.SetExtension(feature, info.ext, res)

		page = pb.Page_builder{
			Features: []*pb.Feature{feature},
		}.Build()
		return false // stop on the first extension
	}
	proto.RangeExtensions(intent, f)
	return
}

type fulfillIntentFunc[Req proto.Message, Res proto.Message] func(msg Req) (Res, error)

func dispatch[Req proto.Message, Res proto.Message](f fulfillIntentFunc[Req, Res]) fulfillFunc {
	return func(msg proto.Message) (proto.Message, error) {
		req, ok := msg.(Req)
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrIntent, msg)
		}
		return f(req)
	}
}
