// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scheduler

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const (
	cpuProbMin = 0
	cpuProbMax = 100
)

// ProcessSpec configures a process.
type ProcessSpec struct {
	numInstructions int
	cpuProb         int // probability [0..100] to get a CPU instruction, else IO
}

// NewProcessSpec creates a [ProcessSpec] from a "n:prob" spec. It returns an
// error if spec is invalid, i.e. the tokens are not numbers or probability is
// out of [0,100] range.
func NewProcessSpec(spec string) (*ProcessSpec, error) {
	before, after, ok := strings.Cut(spec, ":")
	if !ok {
		return nil, fmt.Errorf("spec %q: missing colon separator", spec)
	}
	n, err := strconv.Atoi(before)
	if err != nil {
		return nil, fmt.Errorf("spec %q: %v", spec, err)
	}
	prob, err := strconv.Atoi(after)
	if err != nil {
		return nil, fmt.Errorf("spec %q: %v", spec, err)
	}
	if prob < cpuProbMin || prob > cpuProbMax {
		return nil, fmt.Errorf("spec %q: cpu chance %d: out of range [0..100]", spec, prob)
	}
	return &ProcessSpec{
		numInstructions: n,
		cpuProb:         prob,
	}, nil
}

func (spec *ProcessSpec) String() string {
	return fmt.Sprintf("%d:%d", spec.numInstructions, spec.cpuProb)
}

// Instruction represents a type of instruction.
//
//go:generate stringer -type=Instruction -linecomment
type Instruction int

const (
	_ Instruction = iota
	// keep-sorted start
	InstructionCPU    // cpu
	InstructionIO     // io
	InstructionIODone // iod
	// keep-sorted end
)

// ProcessState identifies the process's state.
//
//go:generate stringer -type=ProcessState -linecomment
type ProcessState int

const (
	_ ProcessState = iota
	// keep-sorted start
	ProcessStateBlocked // blk
	ProcessStateInitial // ini
	ProcessStateReady   // rdy
	ProcessStateRunning // run
	ProcessStateZombie  // ok
	// keep-sorted end
)

// Process represents a process with pid, instructions, etc.
type Process struct {
	pid   int
	state ProcessState

	instructions []Instruction
	pc           int // program counter is the index of the next instruction
}

// NewProcess creates a new process from the spec with pid.
func NewProcess(pid int, spec *ProcessSpec) *Process {
	p := &Process{
		pid:          pid,
		state:        ProcessStateInitial,
		instructions: make([]Instruction, 0, spec.numInstructions), // pre-allocate min
	}
	for i := 0; i < spec.numInstructions; i++ {
		if n := rand.Intn(cpuProbMax); n > spec.cpuProb {
			p.instructions = append(p.instructions, InstructionIO, InstructionIODone)
		} else {
			p.instructions = append(p.instructions, InstructionCPU)
		}
	}
	return p
}

// LastInstruction gives access to the last executed instruction. It returns
// ok=false if there is no last instruction, i.e. the process did not run.
func (p *Process) LastInstruction() (inst Instruction, ok bool) {
	if p.pc == 0 {
		return inst, false
	}
	return p.instructions[p.pc-1], true
}

// PID supplies process identifier.
func (p *Process) PID() int {
	return p.pid
}

// Ready sets the process in ready state.
func (p *Process) Ready() {
	p.state = ProcessStateReady
}

// Step runs a single instruction.
func (p *Process) Step() {
	if p.state == ProcessStateZombie {
		panic(fmt.Errorf("process %d: can not run zombie", p.pid))
	}
	if p.pc == len(p.instructions) {
		p.state = ProcessStateZombie
		return
	}
	switch p.instructions[p.pc] {
	case InstructionCPU:
		p.state = ProcessStateRunning
	case InstructionIO:
		p.state = ProcessStateBlocked
	case InstructionIODone:
		p.state = ProcessStateReady
	}
	p.pc += 1
}

// State returns current state of the process
func (p *Process) State() ProcessState {
	return p.state
}

func (p *Process) String() string {
	return fmt.Sprintf("%d [%s]", p.pid, p.state)
}

// WalkInstructions applies a function f to every instruction in the process.
func (p *Process) WalkInstructions(f func(Instruction)) {
	for _, i := range p.instructions {
		f(i)
	}
}
