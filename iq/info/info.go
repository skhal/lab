// Copyright 2025 Samvel Khalatyan. All rights reserved.

package info

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/skhal/lab/iq/pb"
	"github.com/skhal/lab/iq/registry"
)

// ErrQuestionID represents a group of errors due to invalid question
// identifier.
var ErrQuestionID = errors.New("invalid question id")

// Run prints questions from the registry.
func Run(reg *registry.R, args ...string) error {
	printer := newPrinter(reg)
	if len(args) == 0 {
		return printer.PrintAll()
	}
	ids, err := ParseQuestionIDs(args)
	if err != nil {
		return err
	}
	return printer.PrintSome(ids)
}

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

func (p *printer) PrintAll() error {
	p.reg.Visit(printQuestion)
	return nil
}

type MultiQuestionIDError struct {
	IDs []string
}

func (err *MultiQuestionIDError) Is(e error) bool {
	return e == ErrQuestionID
}

func (err *MultiQuestionIDError) Error() string {
	return fmt.Sprintf("%s: %s", ErrQuestionID, strings.Join(err.IDs, ", "))
}

func (p *printer) PrintSome(ids []registry.QuestionID) error {
	var invalidIDs []string
	for _, id := range ids {
		q := p.reg.Get(id)
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

func printQuestion(q *pb.Question) {
	fmt.Printf("%d\t%s\n", q.GetId(), q.GetDescription())
}
