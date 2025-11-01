// Copyright 2025 Samvel Khalatyan. All rights reserved.

package create

import (
	"errors"
	"flag"
	"fmt"

	"github.com/skhal/lab/go/flags"
	"github.com/skhal/lab/iq/registry"
)

type Config struct {
	Description string
	Tags        flags.StringList
}

func (cfg *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.Description, "description", "", "one-line description")
	fs.Var(&cfg.Tags, "tag", "list of tags")
}

func Run(cfg *Config, reg *registry.R) error {
	if cfg.Description == "" {
		return errors.New("missing description")
	}
	q, err := reg.CreateQuestion(cfg.Description, []string(cfg.Tags))
	if err != nil {
		return err
	}
	fmt.Printf("added question %d\n", q.GetId())
	return nil
}
