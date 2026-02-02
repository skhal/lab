// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nextid implements Protobuf checks for next-id comments.
//
// The next-id comments indicate the field ID to be used next, when a new
// field is added to the message.
//
// Example:
//
//		// Next ID: 2
//	 message Foo {
//			int bar = 1;
//	 }
package nextid
