// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Certval validates certificate dates for expirations and sends notifications.
//
// Synopsis
//
//	certval -c /path/to/certval.txtpb
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

var errFlagConfigFile = errors.New("missing configuration file")

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, ok := err.(FlagError); ok {
			fmt.Fprintln(os.Stderr)
			flag.Usage()
		}
		os.Exit(1)
	}
}

func run() error {
	cfgFile := flag.String("c", "", "configuration file")
	flagValidator := func() error {
		if *cfgFile == "" {
			return errFlagConfigFile
		}
		return nil
	}
	if err := ParseFlags(flagValidator); err != nil {
		return err
	}
	cfg, err := ParseConfig(*cfgFile)
	if err != nil {
		return err
	}
	return validateCertificates(cfg)
}

func validateCertificates(cfg *pb.Config) error {
	notifier := NewSlackNotifier(cfg.GetSlack())
	validator := NewValidator(os.ReadFile, notifier)
	var ee []error
	for _, cert := range cfg.GetCertificate() {
		err := validator.Validate(cert)
		ee = append(ee, err)
	}
	return errors.Join(ee...)
}
