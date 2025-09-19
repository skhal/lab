# NAME

**txtpb-go** - Protobuf textformat with Go


# SYNOPSIS

```go
import "google.golang.org/protobuf/encoding/prototext"

foo := new(pb.Foo)
if err := prototext.Unmarshall(bytes, foo); err != nil { ... }
```


# DESCRIPTION

Ref: https://protobuf.dev/reference/protobuf/textformat-spec/

> [!TIP]
> Store Protobuf text formatted single message in a file with `.txtpb` suffix.

Use `prototext.Unmarshall()` to parse Protobuf text format into a message:


