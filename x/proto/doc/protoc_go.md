# NAME

**protoc-go** - compile Protobuf for Go


# SUMMARY

* Install `protobuf` package with Protobuf compiler `protoc`.
* Install Go Protobuf plugin `protoc-gen-go` with `go install`.
* Compile Protobuf into Go code with `go generate`.


# DESCRIPTION

Ref: https://protobuf.dev/getting-started/gotutorial/

Install Protobuf compiler `protoc`:

```console
% doas pkg install protobuf
```

Install Go Protobuf plugin:

```console
% go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Define Protobuf messages in `x/proto/pb/pingpong.proto`.

> [!TIP]
> Place Protobuf definitions under `pb/` folder.

https://github.com/skhal/lab/blob/443c3225718191cb207097030b4b7950c47db443/x/proto/pb/pingpong.proto#L1-L22

Instead of manually running Protobuf compiler on `.proto` files, we'll use
[`go generate`](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source).
It must be run manually.

Place a `gen.go` file next to the source files with `//go:generate` directives
on how to compile `.proto` files to Go.

> [!TIP]
> Create a command alias with common flags for readability of the file. For
> example, use `//go:generate -comand foo ...` to create an alias `foo`.

https://github.com/skhal/lab/blob/b2f7174d45867695b9fa9799902d3f965e809258/x/proto/pb/gen.go#L5-L6

Run Go generate commands:

```console
% go generate x/proto/pb/....
```

> [!TIP]
> It is common to check in the generated code into Version Control System.
