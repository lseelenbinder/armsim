// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"testing"
	"time"
	"code.google.com/p/go.net/websocket"
)

func TestServer(t *testing.T) {
	// Launch the server
	go Launch()

	// Wait for server to come up (just in case)
	time.Sleep(time.Duration(10) * time.Millisecond)

	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}

	var m Message
	var s Status

	// Test Initial Ping
	m = Message{"hello", ""}
	m.Send(ws)
	websocket.JSON.Receive(ws, &m)
	if m.Cmd != "ello" {
		t.Fatal("Did not acknowledge ping.")
	}

	// Test Load
	m = Message{"load", "../../test/test1.exe"}
	m.Send(ws)
	websocket.JSON.Receive(ws, &s)
}
