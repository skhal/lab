// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pb shows how to define protobuf options (field, message, etc.), use
// it in other messages, proto-compile the code with protoc, and use it in Go.
//
// # WARNING
//
// It is important to run protoc with include paths relative to the repository
// worktree. protoc-gen-go generates init-functions named after the the file
// passed to protoc. The init-function calls initializers for the imports, that
// are named after the import path. All init-names must match else .proto file
// that imports can't initialize the imported files - Go compilation error.
//
// Ref: https://github.com/golang/protobuf/issues/1232
package pb
