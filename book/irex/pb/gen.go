// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pb holds market data schema using Protobuf.
package pb

// protoc compiles Protobuf to Go
//go:generate -command protoc protoc --proto_path=. -I=. -I=../../../ --go_out=../../../ --go_opt=paths=source_relative

// Use paths to files relative to git-worktree:
// https://github.com/golang/protobuf/issues/1232
//
//go:generate protoc book/irex/pb/market.proto
