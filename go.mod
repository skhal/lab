module github.com/skhal/lab

go 1.26.0

require (
	github.com/bufbuild/protocompile v0.14.1
	github.com/gomarkdown/markdown v0.0.0-20260417124207-7d523f7318df
	github.com/google/go-cmp v0.7.0
	github.com/lpar/gzipped/v2 v2.1.0
	github.com/mmcdole/gofeed v1.5.0
	github.com/protocolbuffers/txtpbfmt v0.0.0-20260420112717-c39628bde8b5
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/google/pprof v0.0.0-20260402051712-545e8a4df936 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20250417193237-f615e6bd150b // indirect
	github.com/kevinpollet/nego v0.0.0-20211010160919-a65cd48cee43 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mmcdole/goxpp/v2 v2.0.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260817212433-ac3dfec99bb1 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.6.1 // indirect
	honnef.co/go/tools v0.8.1 // indirect
)

tool (
	github.com/google/pprof
	golang.org/x/tools/cmd/stringer
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
	honnef.co/go/tools/cmd/staticcheck
)
