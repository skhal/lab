// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/skhal/lab/x/fin/internal/pb"
)

// Read reads market data in CSV format from the reader. It returns market data
// in Protobuf format or error. It is all or nothing operation, i.e., an error
// in one of the records invalidaes all data.
func Read(r io.Reader) (*pb.Market, error) {
	csvr := csv.NewReader(r)
	if err := skipHeader(csvr); err != nil {
		return nil, err
	}
	rr, err := readRecords(csvr)
	if err != nil {
		return nil, err
	}
	m := new(pb.Market)
	m.SetRecords(rr)
	return m, nil
}

const headerLines = 8

func skipHeader(r *csv.Reader) error {
	for i := headerLines; i > 0; i -= 1 {
		_, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func readRecords(r *csv.Reader) ([]*pb.Record, error) {
	var records []*pb.Record
	for lineNum := headerLines + 1; ; lineNum += 1 {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}
		rec, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}
		if err := validate(rec, records); err != nil {
			return nil, fmt.Errorf("%d: %v", lineNum, err)
		}
		records = append(records, rec)
	}
	return records, nil
}

func validate(rec *pb.Record, prev []*pb.Record) error {
	if len(prev) < 1 {
		return nil
	}
	isNextMonth := func(prev *pb.Date, next *pb.Date) bool {
		switch next.GetMonth() - prev.GetMonth() {
		case 1: // same year
			return next.GetYear() == prev.GetYear()
		case -11: // next year, Jan - Dec = 1 - 12 = -11
			return next.GetYear() == prev.GetYear()+1
		default:
			return false
		}
	}
	if !isNextMonth(prev[len(prev)-1].GetDate(), rec.GetDate()) {
		return fmt.Errorf("not next month")
	}
	return nil
}
