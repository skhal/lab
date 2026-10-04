// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

const contentTypeJSON = "content-type: application/json"

// ErrNotify means there is an error in sending a notification.
var ErrNotify = errors.New("notification failed")

// SlackNotifier sends notifications to a Slack webook.
type SlackNotifier struct {
	slack *pb.Slack
}

// NewSlackNotifier creates a SlackNotifier.
func NewSlackNotifier(slack *pb.Slack) *SlackNotifier {
	return &SlackNotifier{slack}
}

// SlackData represents JSON data to be sent to Slack.
type SlackData struct {
	// Text is the message body.
	Text string `json:"text"`
}

// Error sends an error notification to Slack.
func (sl *SlackNotifier) Error(err error) error {
	data := SlackData{
		Text: err.Error(),
	}
	body, e := json.Marshal(data)
	if e != nil {
		return e
	}
	if e := send(body, sl.slack.GetWebhook()); e != nil {
		return fmt.Errorf("%w: %s\nServer response: %s", ErrNotify, e, err)
	}
	return nil
}

func send(b []byte, url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	buf := bytes.NewBuffer(b)
	resp, err := client.Post(url, contentTypeJSON, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", bytes.TrimSpace(respBody))
	}
	return nil
}
