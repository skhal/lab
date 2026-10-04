// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "errors"

// DispatchNotifier broadcasts notifications to notifiers.
type DispatchNotifier struct {
	notifiers []notifier
}

// NewDispatchNotifier creates a DispatchNotifier with at least one notifier.
func NewDispatchNotifier(n notifier, nn ...notifier) DispatchNotifier {
	return DispatchNotifier{
		notifiers: append([]notifier{n}, nn...),
	}
}

// Notify dispatches the notification to all notifiers and propagates returned
// errors if any.
func (dn DispatchNotifier) Notify(err error) error {
	var ee []error
	for _, n := range dn.notifiers {
		ee = append(ee, n.Notify(err))
	}
	return errors.Join(ee...)
}
