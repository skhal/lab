// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"

	"github.com/skhal/lab/iq/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

// ErrRegistry is a catch all error in the registry.
var ErrRegistry = errors.New("registry error")

// ErrQuestion indicates an error in a given question.
type ErrQuestion struct {
	Question *pb.Question
}

// Error prints the question information.
func (e *ErrQuestion) Error() string {
	q := e.Question
	if q == nil {
		return fmt.Sprintf("%s: <nil> question", ErrRegistry)
	}
	return fmt.Sprintf("%s: %d: %s", ErrRegistry, q.GetId(), q.GetDescription())
}

func (e *ErrQuestion) Is(err error) bool {
	return err == ErrRegistry
}

// QuestionID is the question unique identifier.
type QuestionID int

// R holds interview questions, keyed by the question ID.
type R struct {
	qset map[QuestionID]*pb.Question
}

// Load reads registry from the input file in Protobuf text format. It returns
// an error if the file does not exist or loading fails to parse text proto.
func Load(file string) (*R, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	questionSet := new(pb.QuestionSet)
	if err := prototext.Unmarshal(b, questionSet); err != nil {
		return nil, err
	}
	return WithQuestions(questionSet.GetQuestions())
}

func WithQuestions(qq []*pb.Question) (*R, error) {
	qset := make(map[QuestionID]*pb.Question)
	for _, q := range qq {
		qid := QuestionID(q.GetId())
		if _, ok := qset[qid]; ok {
			return nil, &ErrQuestion{Question: q}
		}
		qset[qid] = q
	}
	reg := &R{
		qset: qset,
	}
	return reg, nil
}

// Visit passes every question in the registry to the visitor v. The questions
// are ordered by identifiers.
func (r *R) Visit(v func(*pb.Question)) {
	sortedIDs := sortableQuestionIDs(slices.Collect(maps.Keys(r.qset)))
	sort.Sort(sortedIDs)
	for _, qid := range sortedIDs {
		q := r.qset[qid]
		v(q)
	}
}

type sortableQuestionIDs []QuestionID

// Len reports the number of question ids.
func (qq sortableQuestionIDs) Len() int {
	return len(qq)
}

// Swap exchanges two question ids at indices i and j.
func (qq sortableQuestionIDs) Swap(i, j int) {
	qq[i], qq[j] = qq[j], qq[i]
}

// Less reports whether question identifier at index i is less than that at the
// index j.
func (qq sortableQuestionIDs) Less(i, j int) bool {
	return qq[i] < qq[j]
}
