// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

var testRSAPrivateKey *rsa.PrivateKey

func init() {
	const encodedRSAPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsPnoGUOnrpiSqt4XynxA+HRP7S+BSObI6qJ7fQAVSPtRkqso
tWxQYLEYzNEx5ZSHTGypibVsJylvCfuToDTfMul8b/CZjP2Ob0LdpYrNH6l5hvFE
89FU1nZQF15oVLOpUgA7wGiHuEVawrGfey92UE68mOyUVXGweJIVDdxqdMoPvNNU
l86BU02vlBiESxOuox+dWmuVV7vfYZ79Toh/LUK43YvJh+rhv4nKuF7iHjVjBd9s
B6iDjj70HFldzOQ9r8SRI+9NirupPTkF5AKNe6kUhKJ1luB7S27ZkvB3tSTT3P59
3VVJvnzOjaA1z6Cz+4+eRvcysqhrRgFlwI9TEwIDAQABAoIBAEEYiyDP29vCzx/+
dS3LqnI5BjUuJhXUnc6AWX/PCgVAO+8A+gZRgvct7PtZb0sM6P9ZcLrweomlGezI
FrL0/6xQaa8bBr/ve/a8155OgcjFo6fZEw3Dz7ra5fbSiPmu4/b/kvrg+Br1l77J
aun6uUAs1f5B9wW+vbR7tzbT/mxaUeDiBzKpe15GwcvbJtdIVMa2YErtRjc1/5B2
BGVXyvlJv0SIlcIEMsHgnAFOp1ZgQ08aDzvilLq8XVMOahAhP1O2A3X8hKdXPyrx
IVWE9bS9ptTo+eF6eNl+d7htpKGEZHUxinoQpWEBTv+iOoHsVunkEJ3vjLP3lyI/
fY0NQ1ECgYEA3RBXAjgvIys2gfU3keImF8e/TprLge1I2vbWmV2j6rZCg5r/AS0u
pii5CvJ5/T5vfJPNgPBy8B/yRDs+6PJO1GmnlhOkG9JAIPkv0RBZvR0PMBtbp6nT
Y3yo1lwamBVBfY6rc0sLTzosZh2aGoLzrHNMQFMGaauORzBFpY5lU50CgYEAzPHl
u5DI6Xgep1vr8QvCUuEesCOgJg8Yh1UqVoY/SmQh6MYAv1I9bLGwrb3WW/7kqIoD
fj0aQV5buVZI2loMomtU9KY5SFIsPV+JuUpy7/+VE01ZQM5FdY8wiYCQiVZYju9X
Wz5LxMNoz+gT7pwlLCsC4N+R8aoBk404aF1gum8CgYAJ7VTq7Zj4TFV7Soa/T1eE
k9y8a+kdoYk3BASpCHJ29M5R2KEA7YV9wrBklHTz8VzSTFTbKHEQ5W5csAhoL5Fo
qoHzFFi3Qx7MHESQb9qHyolHEMNx6QdsHUn7rlEnaTTyrXh3ifQtD6C0yTmFXUIS
CW9wKApOrnyKJ9nI0HcuZQKBgQCMtoV6e9VGX4AEfpuHvAAnMYQFgeBiYTkBKltQ
XwozhH63uMMomUmtSG87Sz1TmrXadjAhy8gsG6I0pWaN7QgBuFnzQ/HOkwTm+qKw
AsrZt4zeXNwsH7QXHEJCFnCmqw9QzEoZTrNtHJHpNboBuVnYcoueZEJrP8OnUG3r
UjmopwKBgAqB2KYYMUqAOvYcBnEfLDmyZv9BTVNHbR2lKkMYqv5LlvDaBxVfilE0
2riO4p6BaAdvzXjKeRrGNEKoHNBpOSfYCOM16NjL8hIZB1CaV3WbT5oY+jp7Mzd5
7d56RZOE+ERK2uz/7JX9VSsM/LbH9pJibd4e8mikDS9ntciqOH/3
-----END RSA PRIVATE KEY-----
`
	block, _ := pem.Decode([]byte(encodedRSAPrivateKey))
	if block == nil {
		log.Fatalf("failed to decode private key:\n%s", encodedRSAPrivateKey)
	}
	var err error
	testRSAPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %s", err)
	}
}

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name string
		// genFileReader is a fileRederFunc generator that is called from sub-tests
		// and as access to testing.T to fail the test if necessary.
		genFileReader func(*testing.T) readFileFunc
		wantErr       error
	}{
		{
			name: "non-existent certificate",
			genFileReader: func(*testing.T) readFileFunc {
				return func(string) ([]byte, error) {
					return nil, os.ErrNotExist
				}
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "invalid encoding",
			genFileReader: func(*testing.T) readFileFunc {
				cert := `
---BEGIN CERTIFICATE----
ABC12345
-----END CERTIFICATE-----
`
				return func(string) ([]byte, error) {
					return []byte(cert), nil
				}
			},
			wantErr: ErrCertEncoding,
		},
		{
			name: "malformed certificate",
			genFileReader: func(*testing.T) readFileFunc {
				cert := `
-----BEGIN CERTIFICATE-----
ABC12345
-----END CERTIFICATE-----
`
				return func(string) ([]byte, error) {
					return []byte(cert), nil
				}
			},
			wantErr: ErrCertParse,
		},
		{
			name: "about to expire certificate",
			genFileReader: func(t *testing.T) readFileFunc {
				cert := &x509.Certificate{
					NotAfter: time.Now().Add(1 * 24 * time.Hour),
				}
				b := createEncodedCertificate(t, cert, testRSAPrivateKey)
				return func(string) ([]byte, error) {
					return b, nil
				}
			},
			wantErr: ErrCertExpiring,
		},
		{
			name: "expired certificate",
			genFileReader: func(t *testing.T) readFileFunc {
				cert := &x509.Certificate{
					NotAfter: time.Now().Add(-1 * 24 * time.Hour),
				}
				b := createEncodedCertificate(t, cert, testRSAPrivateKey)
				return func(string) ([]byte, error) {
					return b, nil
				}
			},
			wantErr: ErrCertExpired,
		},
		{
			name: "valid certificate",
			genFileReader: func(t *testing.T) readFileFunc {
				cert := &x509.Certificate{
					NotAfter: time.Now().Add((defaultExpireDays + 1) * 24 * time.Hour),
				}
				b := createEncodedCertificate(t, cert, testRSAPrivateKey)
				return func(string) ([]byte, error) {
					return b, nil
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fileReader := tc.genFileReader(t)
			notifier := &testNotifier{}
			validator := NewValidator(fileReader, notifier)
			cert := pb.Certificate_builder{
				Name: new(tc.name),
				Path: new(strings.ReplaceAll(tc.name, " ", "_")),
			}.Build()

			err := validator.Validate(cert)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
		})
	}
}

type testNotifier struct{}

func (tn *testNotifier) Notify(err error) error { return err }

func createEncodedCertificate(t *testing.T, cert *x509.Certificate, key *rsa.PrivateKey) []byte {
	t.Helper()
	b, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create encoded certificate: %s", err)
	}
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: b,
	}
	return pem.EncodeToMemory(block)
}

func TestCertificateError_Unwrap(t *testing.T) {
	cert := pb.Certificate_builder{Name: new("test certificate")}.Build()
	want := errors.New("test error")

	got := errors.Unwrap(NewCertificateError(cert, want))

	if got != want {
		t.Errorf("Unwrap got %s; want %s", got, want)
	}
}
