// Copyright 2025 Samvel Khalatyan. All rights reserved.

package pb_test

import (
	"fmt"

	"github.com/skhal/lab/x/proto/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

func Example_txtpb() {
	txtpb := `
message: "demo message"
`
	foo := new(pb.Foo)
	if err := prototext.Unmarshal([]byte(txtpb), foo); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(foo.GetMessage())
	}
	// Output:
	// demo message
}
