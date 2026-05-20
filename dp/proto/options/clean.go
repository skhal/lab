// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/skhal/lab/dp/proto/options/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Clean removes fields from the message that have field option
// privacy.sensitive set to true.
func Clean(msg proto.Message) {
	m := msg.ProtoReflect()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		if !proto.GetExtension(opts, pb.E_Privacy).(*pb.Privacy).GetSensitive() {
			return true
		}
		m.Clear(fd)
		return true
	})
}
