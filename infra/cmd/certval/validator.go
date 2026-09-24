// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

const defaultExpireDays = 10

type notifier interface {
	Error(error) error
}

type readFileFunc func(string) ([]byte, error)

// Validator validates certificates.
type Validator struct {
	readFile readFileFunc
	notifier notifier
}

// NewValidator creates a validator. It takes a function to read files (faked
// in tests) and a notifier.
func NewValidator(rf readFileFunc, n notifier) *Validator {
	return &Validator{readFile: rf, notifier: n}
}

var (
	// ErrCertEncoding means PEM certificate decode failed.
	ErrCertEncoding = errors.New("certificate encoding error")

	// ErrCertParse means x509 certificated failed to parse.
	ErrCertParse = errors.New("talformed certificate")

	// ErrCertExpiring means the certificate is about to expire, i.e. the
	// NotAfter date is less then certificate's notify.expire_days configuration
	// parameter.
	ErrCertExpiring = errors.New("expiring certificate")

	// ErrCertExpired means the certificate has expired.
	ErrCertExpired = errors.New("expired certificate")
)

// Validate ensures the certificate is accessible and has not expired. It sends
// a notification in case of violations found.
func (vr *Validator) Validate(cert *pb.Certificate) error {
	b, err := vr.readFile(cert.GetPath())
	if err != nil {
		return vr.notifier.Error(NewCertificateError(cert, err))
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return vr.notifier.Error(NewCertificateError(cert, ErrCertEncoding))
	}
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrCertParse, err)
		return vr.notifier.Error(NewCertificateError(cert, err))
	}
	return vr.validateNotAfter(cert, parsedCert)
}

func (vr *Validator) validateNotAfter(cert *pb.Certificate, parsedCert *x509.Certificate) error {
	if parsedCert.NotAfter.Before(time.Now()) {
		return vr.notifier.Error(NewCertificateError(cert, ErrCertExpired))
	}
	var (
		years      int // zero
		months     int // zero
		expireDays = int(cert.GetNotify().GetExpireDays())
	)
	if !cert.GetNotify().HasExpireDays() {
		expireDays = defaultExpireDays
	}
	notifyDate := time.Now().AddDate(years, months, expireDays)
	if parsedCert.NotAfter.Before(notifyDate) {
		var msg string
		if h := time.Until(parsedCert.NotAfter).Hours(); h > 24 {
			msg = fmt.Sprintf("expires in %.0f days", math.Floor(h/24))
		} else {
			msg = fmt.Sprintf("expires in %.0f hours", h)
		}
		err := fmt.Errorf("%w: %s (expire days is %d)", ErrCertExpiring, msg, expireDays)
		return vr.notifier.Error(NewCertificateError(cert, err))
	}
	return nil
}

// CertificateError wraps a certificate related error and prints certificate
// path.
type CertificateError struct {
	cert *pb.Certificate
	err  error
}

// NewCertificateError creates a CertificateError.
func NewCertificateError(cert *pb.Certificate, err error) *CertificateError {
	return &CertificateError{cert: cert, err: err}
}

// Error implements [builtin.error] interface.
func (e *CertificateError) Error() string {
	return fmt.Sprintf("Certificate error: %s: %s", e.cert.GetPath(), e.err)
}

// Unwrap returns wrapped error for [errors.Unwrap].
func (e *CertificateError) Unwrap() error {
	return e.err
}
