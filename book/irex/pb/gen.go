// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pb holds market data schema using Protobuf.
package pb

// protoc compiles Protobuf to Go
//go:generate -command protoc protoc --proto_path=. -I=. -I=../../../ --go_out=../../../ --go_opt=paths=source_relative

//go:generate -command protoc_grpc protoc --proto_path=. -I=../../../ --go_out=../../../ --go_opt=paths=source_relative --go-grpc_out=../../../ --go-grpc_opt=paths=source_relative

// Use paths to files relative to git-worktree:
// https://github.com/golang/protobuf/issues/1232
//
// keep-sorted start
//go:generate protoc book/irex/pb/symbol.proto
//go:generate protoc book/irex/pb/intent.proto
//go:generate protoc book/irex/pb/market.proto
//go:generate protoc book/irex/pb/plot_intent.proto
// keep-sorted end
//
// keep-sorted start
//go:generate protoc_grpc book/irex/pb/market_service.proto
// keep-sorted end
