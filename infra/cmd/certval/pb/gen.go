// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pb defines configuration for the certificate validator.
package pb

// protoc compiles Protobuf to Go.
//go:generate -command protoc protoc --proto_path=. -I=. -I=../../../../ --go_out=../../../../ --go_opt=paths=source_relative

// Use file paths relative to git-worktree:
// https://github.com/golang/protobuf/issues/1232
//
//go:generate protoc infra/cmd/certval/pb/config.proto
