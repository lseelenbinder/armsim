// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	var s Server
	// Launch the server
	go s.Launch(os.Stdout)

	// Wait for server to come up (just in case)
	time.Sleep(time.Duration(100) * time.Millisecond)

	origin := "http://localhost:4567/"
	url := "ws://localhost:4567/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}

	var m Message

	// Test Initial Ping
	m = Message{"hello", ""}
	m.Send(ws)
	websocket.JSON.Receive(ws, &m)
	if m.Type != "status" || m.Content != "ready" {
		t.Fatal("Did not acknowledge ping.")
	}

	// Test Load
	m = Message{"load", "../../test/test1.exe"}
	m.Send(ws)
	websocket.JSON.Receive(ws, &m)

	// Test index.html
	resp, err := http.Get("http://localhost:4567/index.html")
	if resp.Status != "200 OK" {
		t.Fatal("Did not load index.html. " + resp.Status)
	}
}
