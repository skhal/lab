// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"

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
	header []string
	qset   map[QuestionID]*pb.Question

	lastid QuestionID

	updated bool
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
	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return nil, err
	}
	questionSet := new(pb.QuestionSet)
	if err := prototext.Unmarshal(data, questionSet); err != nil {
		return nil, err
	}
	r, err := WithQuestions(questionSet.GetQuestions())
	if err != nil {
		return nil, err
	}
	if header := extractHeader(data); len(header) != 0 {
		r.header = header
	}
	return r, nil
}

func extractHeader(data []byte) []string {
	var header []string
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			break
		}
		header = append(header, line)
	}
	return header
}

// Write stores registry in the file in Protobu text format.
func Write(r *R, cfg *Config) error {
	data, err := marshal(r)
	if err != nil {
		return err
	}
	return write(r.header, data, cfg)
}

func marshal(r *R) ([]byte, error) {
	qset := new(pb.QuestionSet)
	qset.SetQuestions(slices.Collect(r.sortedQuestions()))
	opts := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	return opts.Marshal(qset)
}

func write(header []string, data []byte, cfg *Config) error {
	f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(header) != 0 {
		for _, h := range header {
			fmt.Fprintln(f, h)
		}
		fmt.Fprintln(f)
	}
	_, err = f.Write(data)
	return err
}

func WithQuestions(qq []*pb.Question) (*R, error) {
	return With(QuestionSetOption(qq))
}

// Option customizes registry.
type Option func(*R) error

// QuestionOption adds a single question to the registry.
func QuestionOption(q *pb.Question) Option {
	return func(reg *R) error {
		return reg.add(q)
	}
}

// QuestionSetOption adds multiple questions to the registry.
func QuestionSetOption(qq []*pb.Question) Option {
	return func(reg *R) error {
		for _, q := range qq {
			if err := reg.add(q); err != nil {
				return err
			}
		}
		return nil
	}
}

// HeaderOption adds header to the registry. It returns an error if used
// multiple times to avoid header overwrite.
func HeaderOption(h []string) Option {
	return func(reg *R) error {
		if len(reg.header) != 0 {
			return errors.New("header exists")
		}
		reg.header = h
		return nil
	}
}

// With builds a registry with options.
func With(opts ...Option) (*R, error) {
	r := &R{qset: make(map[QuestionID]*pb.Question)}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Visit passes every question in the registry to the visitor v. The questions
// are ordered by identifiers.
func (r *R) Visit(v func(*pb.Question)) {
	for q := range r.sortedQuestions() {
		v(q)
	}
}

func (r *R) sortedQuestions() iter.Seq[*pb.Question] {
	return func(yield func(*pb.Question) bool) {
		sortedIDs := sortableQuestionIDs(slices.Collect(maps.Keys(r.qset)))
		sort.Sort(sortedIDs)
		for _, id := range sortedIDs {
			q := r.qset[id]
			if !yield(q) {
				break
			}
		}
	}
}

func (r *R) CreateQuestion(desc string, tags []string) (*pb.Question, error) {
	q := new(pb.Question)
	q.SetId(int32(r.lastid + 1))
	q.SetDescription(desc)
	q.SetTags(tags)
	if err := r.add(q); err != nil {
		return nil, err
	}
	r.updated = true
	return q, nil
}

func (r *R) add(q *pb.Question) error {
	id := QuestionID(q.GetId())
	if got, ok := r.qset[id]; ok {
		return &ErrDuplicateQuestion{Has: got, New: q}
	}
	r.qset[id] = q
	if id > r.lastid {
		r.lastid = id
	}
	return nil
}

// Get retrieves the question with a given identifier form the registry. It
// returns nil if the question does not exist.
func (r *R) Get(qid QuestionID) *pb.Question {
	return r.qset[qid]
}

func (r *R) Updated() bool {
	return r.updated
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
