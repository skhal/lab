// Copyright 2025 Samvel Khalatyan. All rights reserved.

package pb

//go:generate -command protoc_cmd protoc --proto_path=. -I=../../ --go_out=. --go_opt=paths=source_relative
//go:generate protoc_cmd question.proto
