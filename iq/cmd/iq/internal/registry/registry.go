// Copyright 2025 Samvel Khalatyan. All rights reserved.

package registry

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"sort"

	"github.com/protocolbuffers/txtpbfmt/parser"
	"github.com/skhal/lab/iq/cmd/iq/internal/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

// ErrRegistry is a catch all error in the registry.
var ErrRegistry = errors.New("registry error")

// DuplicateQuestionError indicates an error in a given question.
type DuplicateQuestionError struct {
	Has *pb.Question
	New *pb.Question
}

// Error prints the question information.
func (e *DuplicateQuestionError) Error() string {
	qhas := e.Has
	qnew := e.New
	return fmt.Sprintf("%s: duplicate question %d: has %q, new %q", ErrRegistry, qhas.GetId(), qhas.GetDescription(), qnew.GetDescription())
}

func (e *DuplicateQuestionError) Is(err error) bool {
	return err == ErrRegistry
}

// Config holds registry configuration paramters, to be extracted from flags.
type Config struct {
	// File is the registry filename
	File string
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.File, "f", "iq/registry/questions.txtpb", "questions list")
}

// QuestionID is the question unique identifier.
type QuestionID int

// Tag is the question tag.
type Tag string

type index struct {
	byid  map[QuestionID]*pb.Question
	bytag map[Tag][]*pb.Question
}

// R holds interview questions, keyed by the question ID.
type R struct {
	rootPath string

	header []byte
	index  index

	lastid  QuestionID
	updated bool
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

func RootPathOption(p string) Option {
	return func(reg *R) error {
		if reg.rootPath != "" {
			return errors.New("root path exists")
		}
		reg.rootPath = p
		return nil
	}
}

// HeaderOption adds header to the registry. It returns an error if used
// multiple times to avoid header overwrite.
func HeaderOption(data []byte) Option {
	return func(reg *R) error {
		if len(reg.header) != 0 {
			return errors.New("header exists")
		}
		data = bytes.TrimSpace(data)
		reg.header = data
		return nil
	}
}

// With builds a registry with options.
func With(opts ...Option) (*R, error) {
	r := &R{
		index: index{
			byid:  make(map[QuestionID]*pb.Question),
			bytag: make(map[Tag][]*pb.Question),
		},
	}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Load reads registry from the input file in Protobuf text format. It returns
// an error if the file does not exist or loading fails to parse text proto.
func Load(cfg *Config) (*R, error) {
	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return nil, err
	}
	qset := new(pb.QuestionSet)
	if err := prototext.Unmarshal(data, qset); err != nil {
		return nil, err
	}
	opts := []Option{
		QuestionSetOption(qset.GetQuestions()),
	}
	if qset.HasRootPath() {
		opts = append(opts, RootPathOption(qset.GetRootPath()))
	}
	if header := extractHeader(data); len(header) != 0 {
		opts = append(opts, HeaderOption(header))
	}
	return With(opts...)
}

func extractHeader(data []byte) []byte {
	var prefix = []byte("#")
	size := 0
	for line := range bytes.Lines(data) {
		if !bytes.HasPrefix(line, prefix) {
			break
		}
		size += len(line)
	}
	return bytes.TrimSpace(data[:size])
}

// Write stores registry in the file in Protobu text format.
func Write(r *R, cfg *Config) error {
	data, err := marshal(r)
	if err != nil {
		return err
	}
	data, err = parser.FormatWithConfig(data, parser.Config{
		SkipAllColons: true,
	})
	if err != nil {
		return errors.Join(ErrRegistry, err)
	}
	return write(data, cfg)
}

func (r *R) questionSet() *pb.QuestionSet {
	qset := new(pb.QuestionSet)
	qset.SetQuestions(slices.Collect(r.sortedQuestions()))
	if r.rootPath != "" {
		qset.SetRootPath(r.rootPath)
	}
	return qset
}

func marshal(r *R) ([]byte, error) {
	qset := r.questionSet()
	opts := prototext.MarshalOptions{
		Multiline: true,
	}
	data, err := opts.Marshal(qset)
	if err != nil {
		return nil, err
	}
	const eol = '\n'
	buf := new(bytes.Buffer)
	if len(r.header) != 0 {
		buf.Write(r.header)
		buf.WriteByte(eol) // end of header
		buf.WriteByte(eol) // header / body separator
	}
	buf.Write(data)
	return buf.Bytes(), nil
}

func write(data []byte, cfg *Config) error {
	f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
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
		sortedIDs := sortableQuestionIDs(slices.Collect(maps.Keys(r.index.byid)))
		sort.Sort(sortedIDs)
		for _, id := range sortedIDs {
			q := r.index.byid[id]
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
	if err := r.addToIndexByID(q); err != nil {
		return err
	}
	r.addToIndexByTag(q)
	return nil
}

func (r *R) addToIndexByID(q *pb.Question) error {
	id := QuestionID(q.GetId())
	if got, ok := r.index.byid[id]; ok {
		return &DuplicateQuestionError{Has: got, New: q}
	}
	r.index.byid[id] = q
	if id > r.lastid {
		r.lastid = id
	}
	return nil
}

func (r *R) addToIndexByTag(q *pb.Question) {
	for _, tag := range q.GetTags() {
		qq := r.index.bytag[Tag(tag)]
		qq = append(qq, q)
		r.index.bytag[Tag(tag)] = qq
	}
}

// GetByID retrieves the question with a given identifier form the registry. It
// returns nil if the question does not exist.
func (r *R) GetByID(qid QuestionID) *pb.Question {
	return r.index.byid[qid]
}

// GetByTag retrieves questions for a given tag or nil of no questions are
// registered for the tag.
func (r *R) GetByTag(t Tag) []*pb.Question {
	return r.index.bytag[t]
}

// GetTags retrieves a set of available tags.
func (r *R) GetTags() []Tag {
	return slices.Collect(maps.Keys(r.index.bytag))
}

func (r *R) RootPath() string {
	return r.rootPath
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
