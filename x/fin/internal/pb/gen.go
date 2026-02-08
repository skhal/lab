// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pb holds Protobuf schema for Shiller market data.
package pb

// protoc_cmd_go compiles Protobuf to Go
//go:generate -command protoc_cmd_go protoc --proto_path=. -I=../../../../ --go_out=. --go_opt=paths=source_relative

// protoc_cmd_pb generates Protobuf descriptors
//go:generate -command protoc_cmd_pb protoc --proto_path=. -I=../../../../ --include_imports

//go:generate protoc_cmd_go market.proto
//go:generate protoc_cmd_pb -o market.pb market.proto
