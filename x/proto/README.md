# NAME

**proto** - protobuf examples


# DESCRIPTION

This package contains examples of working with Protocol Buffers.

## Example 1: Build protobuf with Go

Ref: https://protobuf.dev/getting-started/gotutorial/

Install Protobuf compiler `protoc`:

```console
% doas pkg install protobuf
```

Install Go Protobuf plugin:

```console
% go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Define Protobuf messages in `path/pb/foo.proto`. It is common to place Protobuf
definitions under `pb/` folder (see [`x/proto/pb`](./pb)).

We'll use [`go generate`](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source)
to generate Go files by processing source.

Create a `pb/gen.go` file with `//go:generate` directives to compile `.proto`
files to `.go`. It is helpful to create an alias to a command with common flags
with `//go:generate -command foo ...` to create an alias `foo`.

https://github.com/skhal/lab/blob/b2f7174d45867695b9fa9799902d3f965e809258/x/proto/pb/gen.go#L5-L6

Generate go code:

```console
% go generate ....
```
