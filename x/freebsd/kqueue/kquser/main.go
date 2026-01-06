// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//go:build !linux

/*
User demonstrates kqueue(2) send and receive a user event.

Synopsis:

	kquser

Background:

BSD-like systems (*BSD, MacOS, etc.) provide kqueue(2) to monitor events from a
wide range of systems: file system, processes, signals, jails, timer, user, etc.
(MacOS does not include jails).

For the record, Linux has similar mechanism, implemented via select(2), poll(2),
epoll(7), or inotify(7).

Kqueue(2) is a stateful mechanism: Kernel keeps track of subscriptions. It is
important to close kqueue(2) when done to clean up resources.

Kqueue(2) provides a single function to manage and poll events. It guarantees to
process manage-events first, and then update the poll-events. Even though it is
possible to use both kinds of events in the same call to kevent(2), it is best
to separate the two.

Example:

This example demonstrates the following aspects of working with kqueue(2):

- Create and cleanup kqueue(2) Kernel resources.
- Subscribe to User events.
- Send and receive user events.

Output:

	2026/01/05 16:13:50 sent event event-id-two
	2026/01/05 16:13:50 received event event-id-two
*/
package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	kq, err := NewKQueue()
	if err != nil {
		log.Fatal(err)
	}
	defer kq.Close()
	if err := run(kq); err != nil {
		log.Fatal(err)
	}
}

func run(kq *KQueue) error {
	events := []EventID{EventIDOne, EventIDTwo}
	if err := register(kq, events); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(kq *KQueue, eid EventID) {
		defer wg.Done()
		send(kq, eid)
	}(kq, MustGetRandomItem(events))
	wg.Add(1)
	go func(kq *KQueue) {
		defer wg.Done()
		receive(kq)
	}(kq)
	wg.Wait()
	return nil
}

func register(kq *KQueue, events []EventID) error {
	for _, eid := range events {
		if err := kq.Register(eid); err != nil {
			return err
		}
	}
	return nil
}

func send(kq *KQueue, eid EventID) {
	if err := kq.Send(eid); err != nil {
		log.Println(err)
		return
	}
	log.Println("sent event", eid)
}

func receive(kq *KQueue) {
	// Poll at most N times to prevent deadlock if send fails. Avoid context in demo.
	const maxPolls = 5
	for n := 0; n < maxPolls; n++ {
		timeout := time.Duration((5 + rand.Intn(10))) * time.Millisecond
		event, ok := kq.Poll(timeout)
		if !ok {
			continue
		}
		log.Println("received event", event)
		break
	}
}

func MustGetRandomItem[T any](items []T) T {
	if len(items) == 0 {
		panic(errors.New("empty slice"))
	}
	return items[rand.Intn(len(items))]
}
