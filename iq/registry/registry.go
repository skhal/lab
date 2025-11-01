// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry

import (
	"errors"
	"flag"
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

// ErrDuplicateQuestion indicates an error in a given question.
type ErrDuplicateQuestion struct {
	Has *pb.Question
	New *pb.Question
}

// Error prints the question information.
func (e *ErrDuplicateQuestion) Error() string {
	qhas := e.Has
	qnew := e.New
	return fmt.Sprintf("%s: duplicate question %d: has %q, new %q", ErrRegistry, qhas.GetId(), qhas.GetDescription(), qnew.GetDescription())
}

func (e *ErrDuplicateQuestion) Is(err error) bool {
	return err == ErrRegistry
}

// QuestionID is the question unique identifier.
type QuestionID int

// R holds interview questions, keyed by the question ID.
type R struct {
	qset map[QuestionID]*pb.Question
}

// Config holds registry configuration paramters, to be extracted from flags.
type Config struct {
	// File is the registry filename
	File string
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.File, "file", "iq/registry/questions.txtpb", "questions list (txtpb)")
}

// Load reads registry from the input file in Protobuf text format. It returns
// an error if the file does not exist or loading fails to parse text proto.
func Load(cfg *Config) (*R, error) {
	b, err := os.ReadFile(cfg.File)
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
	r := &R{qset: make(map[QuestionID]*pb.Question)}
	for _, q := range qq {
		if err := r.add(q); err != nil {
			return nil, err
		}
	}
	return r, nil
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

func (r *R) add(q *pb.Question) error {
	id := QuestionID(q.GetId())
	if got, ok := r.qset[id]; ok {
		return &ErrDuplicateQuestion{Has: got, New: q}
	}
	r.qset[id] = q
	return nil
}

// Get retrieves the question with a given identifier form the registry. It
// returns nil if the question does not exist.
func (r *R) Get(qid QuestionID) *pb.Question {
	return r.qset[qid]
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
