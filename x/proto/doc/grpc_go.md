# NAME

**grpc-go** - gRPC with Go


# SYNOPSIS

```console
% go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
% go generate ./x/proto/pb/...
```

# DESCRIPTION

Ref: https://grpc.io/docs/languages/go/basics/

It is assumed the environment is setup with Protobuf compiler and Go plugin.

Install Protobuf plugin for Go gRPC using Go modules:

```console
% go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Define the service `x/proto/pb/fooer.proto`.

Add Go generate directives to compile the service definition into Go.

Run Go generate:

```
% go generate ./x/proto/pb/...
```
