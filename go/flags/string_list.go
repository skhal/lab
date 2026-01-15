// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags

import (
	"fmt"
	"strings"
)

const separator = ","

// StringList implements a list of strings flag. It accumulates non-empty values
// after trimming spaces from every use of the flag. Use comma to separate
// multiple values in a single flag.
type StringList []string

// Set implemnets flag.Value interface.
func (f *StringList) Set(value string) error {
	tokens := strings.Split(value, separator)
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		*f = append(*f, token)
	}
	return nil
}

// Get implements flag.Getter interface.
func (f *StringList) Get() any {
	return []string(*f)
}

// String implemnets flag.Value interface.
func (f *StringList) String() string {
	return fmt.Sprint([]string(*f))
}
