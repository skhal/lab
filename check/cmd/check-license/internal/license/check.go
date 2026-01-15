// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package license

import (
	"bytes"
	"errors"
	"regexp"
)

var lineRx = regexp.MustCompile(`^([\t ]*(/[/\*]|[#"])?) Copyright`)

const eol = '\n'

func Check(buf []byte) (err error) {
	ln := 1
	defer func() {
		if errors.Is(err, ErrInvalid) {
			err = NewInvalidError(ln)
		}
	}()
	for len(buf) > 0 {
		if ok, err := match(buf); ok {
			return nil
		} else if err != nil {
			return err
		}
		idx := bytes.IndexByte(buf, eol)
		if idx == -1 {
			break
		}
		idx += 1 // skip eol
		buf = buf[idx:]
		ln += 1
	}
	return ErrNotFound
}

func match(buf []byte) (ok bool, err error) {
	matches := lineRx.FindSubmatch(buf)
	if matches == nil {
		return
	}
	blockRx, err := compileBlockRx(matches[1])
	if err != nil {
		return
	}
	if !blockRx.Match(buf) {
		return false, ErrInvalid
	}
	return true, nil
}

func compileBlockRx(prefix []byte) (*regexp.Regexp, error) {
	b, err := genLicenseBlock(LicenseData{
		Year:   `\d{4}`,
		Holder: `\w+( \w+)?`,
	})
	if err != nil {
		return nil, err
	}
	b = addPrefix(b, prefix)
	return regexp.Compile("^" + string(b)) // must match the beginning
}

func addPrefix(block []byte, prefix []byte) []byte {
	var b bytes.Buffer
	for line := range bytes.Lines(block) {
		if isEmptyLine(line) {
			if len(bytes.TrimSpace(prefix)) != 0 {
				b.Write(prefix)
			}
		} else {
			b.Write(prefix)
			b.WriteByte(' ')
		}
		b.Write(line)
	}
	return b.Bytes()
}

func isEmptyLine(b []byte) bool {
	if len(b) > 1 {
		return false
	}
	return len(b) == 0 || (len(b) == 1 && b[0] == eol)
}
