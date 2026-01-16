// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/skhal/lab/go/flags"
	"github.com/skhal/lab/iq/cmd/iq/internal/registry"
)

type Config struct {
	Description string
	Tags        flags.StringList
}

func (cfg *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.Description, "d", "", "one-line description")
	fs.Var(&cfg.Tags, "t", "list of tags")
}

func Run(cfg *Config, reg *registry.R) error {
	if cfg.Description == "" {
		return errors.New("missing description")
	}
	q, err := reg.CreateQuestion(cfg.Description, []string(cfg.Tags))
	if err != nil {
		return err
	}
	p, err := createQuestionPath(reg.RootPath(), int(q.GetId()))
	if err != nil {
		return err
	}
	fmt.Printf("iniailized path %s\n", p)
	return nil
}

func createQuestionPath(prefix string, id int) (string, error) {
	p := filepath.Join(prefix, strconv.Itoa(int(id)))
	if err := os.Mkdir(p, 0755); err != nil {
		return "", err
	}
	return p, nil
}
