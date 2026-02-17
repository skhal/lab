// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Alpine fetches latest releases and prints a URL to ISO image for requested
// flavor, defaulted to "alpine-virt".
//
// Synopsis:
//
//	alpine [-f flavor]
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const urlBase = "https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/x86_64"

const defaultFlavor = "alpine-virt"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mustUrlJoinPath(base string, elem string) string {
	p, err := url.JoinPath(base, elem)
	if err != nil {
		panic(err)
	}
	return p
}

func run() error {
	flavor := flagParse()
	releases, err := getReleases()
	if err != nil {
		return err
	}
	if flavor == "" {
		printReleases(releases)
		return nil
	}
	return printFlavor(releases, flavor)
}

func flagParse() string {
	flavor := flag.String("flavor", defaultFlavor, "alpine flavor to choose")
	flag.Parse()
	return *flavor
}

func getReleases() ([]*release, error) {
	res, err := http.Get(mustUrlJoinPath(urlBase, "latest-releases.yaml"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var releases []*release
	if err := yaml.Unmarshal(b, &releases); err != nil {
		return nil, err
	}
	return releases, nil
}

func printReleases(rr []*release) {
	for _, rel := range rr {
		fmt.Println(rel)
	}
}

func printFlavor(rr []*release, flavor string) error {
	var flavors []string
	for _, rel := range rr {
		flavors = append(flavors, rel.Flavor)
		if rel.Flavor == flavor {
			fmt.Println(rel)
			return nil
		}
	}
	opts := strings.Join(flavors, ",")
	return fmt.Errorf("missing flavor %s\noptions: %s", flavor, opts)
}

type release struct {
	Version string
	Flavor  string
	ISO     string
	Sha256  string
}

// String implements [fmt.Stringer] interface.
func (r *release) String() string {
	return mustUrlJoinPath(urlBase, r.ISO)
}
