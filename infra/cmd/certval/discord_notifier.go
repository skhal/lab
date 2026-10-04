// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/skhal/lab/infra/cmd/certval/pb"
)

// DiscordNotifier sends notifications to Discord.
type DiscordNotifier struct {
	discord *pb.Discord
}

// NewDiscordNotifier creates a DiscrodNotifier with Discord configuration.
func NewDiscordNotifier(dc *pb.Discord) *DiscordNotifier {
	return &DiscordNotifier{discord: dc}
}

// DiscordData is the Discord JSON params structure for the payload. See:
// https://docs.discord.com/developers/resources/webhook#execute-webhook
type DiscordData struct {
	// Content is the message body.
	Content string `json:"content"`
}

// Notify sends an error notification to Discord.
func (dn *DiscordNotifier) Notify(err error) error {
	data := DiscordData{
		Content: err.Error(),
	}
	body, e := json.Marshal(data)
	if e != nil {
		return e
	}
	if e := send(body, dn.discord.GetWebhook()); e != nil {
		return fmt.Errorf("%w: %s\nServer response: %s", ErrNotify, e, err)
	}
	return nil
}
