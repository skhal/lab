// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/skhal/lab/iq/cmd/iq/internal/pb"
	"github.com/skhal/lab/iq/cmd/iq/internal/registry"
)

// ErrQuestionID represents a group of errors due to invalid question
// identifier.
var ErrQuestionID = errors.New("invalid question id")
var ErrConfig = errors.New("invalid config")

// Config holds parameters to for info command.
type Config struct {
	Tag  string
	Tags bool
}

func (cfg *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.Tag, "t", "", "tag")
	fs.BoolVar(&cfg.Tags, "tt", false, "list tags")
}

func (cfg *Config) Validate() error {
	if cfg.Tags {
		if cfg.Tag != "" {
			return fmt.Errorf("%v: -t and -tt flags are exclusive", ErrConfig)
		}
	}
	return nil
}

// Run prints questions from the registry.
func Run(cfg *Config, reg *registry.R, args ...string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	ids, err := ParseQuestionIDs(args)
	if err != nil {
		return err
	}
	printer := newPrinter(reg)
	switch {
	case cfg.Tags:
		err = printer.PrintTags()
	case cfg.Tag != "":
		err = printer.PrintByTag(registry.Tag(cfg.Tag), ids)
	case len(ids) != 0:
		err = printer.PrintByID(ids)
	default:
		err = printer.PrintAll()
	}
	return err
}

// ParseQuestionIDs parses a list of questions ID strings as integers.
func ParseQuestionIDs(strstr []string) ([]registry.QuestionID, error) {
	ids := make([]registry.QuestionID, 0, len(strstr))
	for _, str := range strstr {
		id, err := strconv.Atoi(str)
		if err != nil {
			return nil, &QuestionIDError{ID: str}
		}
		ids = append(ids, registry.QuestionID(id))
	}
	return ids, nil
}

// QuestionIDError is ErrQuestionID with invalid question identifier.
type QuestionIDError struct {
	ID string
}

// Is make QuestionIDError equivalent to ErrEquestID.
func (err *QuestionIDError) Is(e error) bool {
	return e == ErrQuestionID
}

// Error implements error interface.
func (err *QuestionIDError) Error() string {
	return fmt.Sprintf("%s: %s", ErrQuestionID, err.ID)
}

type printer struct {
	reg *registry.R
}

func newPrinter(reg *registry.R) *printer {
	return &printer{
		reg: reg,
	}
}

// PrintAll prints all questions in the registry.
func (p *printer) PrintAll() error {
	p.reg.Visit(printQuestion)
	return nil
}

// MultiQuestionIDError holds invalid question IDs.
type MultiQuestionIDError struct {
	IDs []string
}

func (err *MultiQuestionIDError) Is(e error) bool {
	return e == ErrQuestionID
}

// Error implemnets errors interface.
func (err *MultiQuestionIDError) Error() string {
	return fmt.Sprintf("%s: %s", ErrQuestionID, strings.Join(err.IDs, ", "))
}

// PrintByID prints questions for selected ids.
func (p *printer) PrintByID(ids []registry.QuestionID) error {
	var invalidIDs []string
	for _, id := range ids {
		q := p.reg.GetByID(id)
		if q == nil {
			invalidIDs = append(invalidIDs, strconv.Itoa(int(id)))
			continue
		}
		printQuestion(q)
	}
	if len(invalidIDs) != 0 {
		return &MultiQuestionIDError{IDs: invalidIDs}
	}
	return nil
}

func (p *printer) PrintTags() error {
	tt := append([]registry.Tag(nil), p.reg.GetTags()...)
	if len(tt) == 0 {
		return nil
	}
	sort.Slice(tt, func(i, j int) bool {
		return strings.Compare(string(tt[i]), string(tt[j])) < 0
	})
	buf := new(bytes.Buffer)
	for _, t := range tt {
		fmt.Fprintln(buf, t)
	}
	fmt.Print(buf.String())
	return nil
}

// PrintByTag prints questions for a given tan and optionally selected by ids.
func (p *printer) PrintByTag(tag registry.Tag, ids []registry.QuestionID) error {
	if len(ids) == 0 {
		return p.printByTagAll(tag)
	}
	qset := make(map[registry.QuestionID]*pb.Question)
	for _, q := range p.reg.GetByTag(tag) {
		qset[registry.QuestionID(q.GetId())] = q
	}
	var invalidIDs []string
	for _, id := range ids {
		q, ok := qset[id]
		if !ok {
			invalidIDs = append(invalidIDs, strconv.Itoa(int(id)))
			continue
		}
		printQuestion(q)
	}
	if len(invalidIDs) != 0 {
		return &MultiQuestionIDError{IDs: invalidIDs}
	}
	return nil
}

func (p *printer) printByTagAll(tag registry.Tag) error {
	for _, q := range p.reg.GetByTag(tag) {
		printQuestion(q)
	}
	return nil
}

func printQuestion(q *pb.Question) {
	fmt.Printf("%d\t%s\n", q.GetId(), q.GetDescription())
}
