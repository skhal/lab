// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//go:build !linux

package main

import (
	"fmt"
	"syscall"
	"time"
)

// KQueue wraps kqueue(2) descriptor.
type KQueue struct {
	fd int
}

// NewKQueue initializes a kqueue in the Kernel. Kqueue is stateful. Make sure
// to clean up Kernel resources when done with the queue.
func NewKQueue() (*KQueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	return &KQueue{
		fd: fd,
	}, nil
}

// Close releases Kernel resources associated with this instance of kqueue(2).
func (kq *KQueue) Close() error {
	if err := syscall.Close(kq.fd); err != nil {
		return err
	}
	kq.fd = 0
	return nil
}

// Register configures kqueue(2) to monitor user events with a given EventID.
func (kq *KQueue) Register(eid EventID) error {
	changes := []syscall.Kevent_t{
		{
			Ident:  uint64(eid),
			Filter: syscall.EVFILT_USER,
			Flags:  syscall.EV_ADD | syscall.EV_CLEAR, // event flags
			// user events may have any value set in the lower 24-bits of fflags.
			// kqueue(2) propagates these to the receiver.
			Fflags: syscall.NOTE_FFNOP, // ignore input fflags
		},
	}
	events := []syscall.Kevent_t(nil) // do not poll events
	var timeout syscall.Timespec      // zero-value, return immediately
	if n, err := syscall.Kevent(kq.fd, changes, events, &timeout); err != nil {
		return err
	} else if n != 0 {
		return fmt.Errorf("unexpected %d pending events for user identifier: %d", n, eid)
	}
	return nil
}

type EventID uint64

const (
	_ EventID = iota
	EventIDOne
	EventIDTwo
)

func (eid EventID) String() string {
	switch eid {
	default:
		return fmt.Sprintf("event-id-%d", eid)
	case EventIDOne:
		return "event-id-one"
	case EventIDTwo:
		return "event-id-two"
	}
}

// Send triggers a user event with EventID.  It returns an error if kevent(2)
// fails.
func (kq *KQueue) Send(eid EventID) error {
	changes := []syscall.Kevent_t{
		{
			Ident:  uint64(eid),          // user event can have any uint64-id
			Filter: syscall.EVFILT_USER,  // it will be a user event
			Fflags: syscall.NOTE_TRIGGER, // send the event
		},
	}
	events := []syscall.Kevent_t(nil) // do not get events
	timeout := syscall.NsecToTimespec((10 * time.Millisecond).Nanoseconds())
	if n, err := syscall.Kevent(kq.fd, changes, events, &timeout); err != nil {
		return err
	} else if n != 0 {
		return fmt.Errorf("unexpected %d pending events in send user id: %d", n, eid)
	}
	return nil
}

// Poll check whether kqueue(2) has a user event available. It is a blocking
// call. Zero duration returns immediately.
func (kq *KQueue) Poll(timeout time.Duration) (*UserEvent, bool) {
	changes := []syscall.Kevent_t(nil)                    // do not change kqueue
	events := make([]syscall.Kevent_t, 1)                 // get at most 1 event
	tout := syscall.NsecToTimespec(timeout.Nanoseconds()) // 0 return immediately
	if n, err := syscall.Kevent(kq.fd, changes, events, &tout); err != nil {
		panic(err)
	} else if n == 0 {
		return nil, false
	}
	return &UserEvent{
		ID: EventID(events[0].Ident),
	}, true
}

// UserEvent wraps the event identifier.
type UserEvent struct {
	ID EventID
}

// String implements fmt.Stringer interface.
func (evt *UserEvent) String() string {
	return fmt.Sprintf("%s", evt.ID)
}
