// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
