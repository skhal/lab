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

https://github.com/skhal/lab/blob/2c99ba3d642faac93e5778c39149da4393f85d23/x/proto/pb/fooer.proto#L12-L24

Add Go generate directives to compile the service definition into Go.

https://github.com/skhal/lab/blob/2c99ba3d642faac93e5778c39149da4393f85d23/x/proto/pb/gen.go#L8-L9

Run Go generate:

```
% go generate ./x/proto/pb/...
```
