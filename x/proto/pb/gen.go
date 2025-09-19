// Copyright 2025 Samvel Khalatyan. All rights reserved.

package pb

//go:generate -command protoc_cmd protoc --proto_path=. -I=../../ --go_out=. --go_opt=paths=source_relative
//go:generate protoc_cmd foo.proto

//go:generate -command protoc_grpc_cmd protoc --proto_path=. -I=../../ --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
//go:generate protoc_grpc_cmd fooer.proto
