// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pb

//go:generate -command protoc_cmd protoc --proto_path=. -I=../../../../../ --go_out=. --go_opt=paths=source_relative
//go:generate protoc_cmd question.proto
