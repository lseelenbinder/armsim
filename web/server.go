package web

import (
	"encoding/json"
	"fmt"
	"github.com/lseelenbinder/armsim/armsim"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Generic Message
type Message struct {
	Type    string
	Content string
}

func (m *Message) Send(ws *websocket.Conn) {
	websocket.JSON.Send(ws, m)
}

type Server struct {
	Computer *armsim.Computer
	FilePath string
	Halt     chan bool
	Finished chan bool
	Log      *log.Logger
	Keyboard chan byte
	Console  chan byte
}

var globalServer Server

func (s *Server) Serve(ws *websocket.Conn) {
	go s.SendConsoleOutput(ws)
	for {
		var m Message

		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			s.Log.Println(err)
			ws.Close()
			return
		}

		s.Log.Println(m)

		switch m.Type {
		case "hello": // Acknowledge ping
			s.SayHi(ws)
		case "load": // Load an ELF file by pathname
			s.Load(m, ws)
		case "reset": // Reset the simulator
			s.Reset(ws)
		case "start": // Run the program
			go s.Start(ws)
		case "status":
			s.UpdateStatus(ws)
		case "step": // Step the program
			s.Computer.Step()
			s.UpdateStatus(ws)
		case "stop": // Stop the program while running
			s.Stop(ws)
		case "trace": // Enable/Disable tracing
			s.Trace(m, ws)
		case "system-trace": // Enable/Disable system tracing
			s.SystemTrace(m, ws)
		case "input":
			s.Input(m, ws)
		case "quit": // Quit connection
			ws.Close()
			break
		default:
			m = Message{"error", "no command"}
			m.Send(ws)
		}
	}
}

func (s *Server) SayHi(ws *websocket.Conn) {
	m := Message{"status", "ready"}
	m.Send(ws)
	if s.FilePath != "" {
		m := Message{"status", "loaded"}
		m.Send(ws)
		m = Message{"filename", s.FilePath}
		m.Send(ws)
	}
	s.UpdateStatus(ws)
}

func (s *Server) Load(m Message, ws *websocket.Conn) {
	path := m.Content
	s.Log.Println(path)
	s.FilePath = path

	s.Computer.Reset()
	err := s.Computer.LoadELF(path)
	if err != nil {
		m := Message{"error", fmt.Sprintf("Unable to load %s. Please check your path.", s.FilePath)}
		m.Send(ws)
	}

	m = Message{"status", "loaded"}
	m.Send(ws)
	s.UpdateStatus(ws)
}

func (s *Server) Reset(ws *websocket.Conn) {
	s.Computer.Reset()
	err := s.Computer.LoadELF(s.FilePath)

	if err != nil {
		m := Message{"error", fmt.Sprintf("Unable to load %s. Please check your path.", s.FilePath)}
		m.Send(ws)
	}

	s.UpdateStatus(ws)
}

func (s *Server) Start(ws *websocket.Conn) {
	go s.Computer.Run(s.Halt, s.Finished)
	m := Message{"status", "running"}
	m.Send(ws)

	// Wait for completion
	<-s.Finished
	s.UpdateStatus(ws)
	m = Message{"status", "finished"}
	m.Send(ws)
}

func (s *Server) Stop(ws *websocket.Conn) {
	s.Halt <- true
	m := Message{"status", "stopped"}
	m.Send(ws)
	s.UpdateStatus(ws)
}

func (s *Server) Trace(m Message, ws *websocket.Conn) {
	// This has the potential to create a race condition. However, I don't think
	// it would even matter (worst case is the last trace is cut off).

	if m.Content == "on" {
		s.Computer.EnableTracing()
	} else {
		s.Computer.DisableTracing()
	}
}

func (s *Server) SystemTrace(m Message, ws *websocket.Conn) {
	// This has the potential to create a race condition. However, I don't think
	// it would even matter (worst case is the last trace is cut off).

	if m.Content == "on" {
		s.Computer.EnableSystemTracing()
	} else {
		s.Computer.DisableSystemTracing()
	}
}

func (s *Server) Quit(ws *websocket.Conn) {
}

func (s *Server) UpdateStatus(ws *websocket.Conn) {
	out, _ := json.Marshal(s.Computer.Status())
	m := Message{"update", string(out)}
	m.Send(ws)
}

func (s *Server) SendConsoleOutput(ws *websocket.Conn) {
	var b byte
	for {
		b = <-s.Console
		m := Message{"output", string(b)}
		m.Send(ws)
	}
}

func (s *Server) Input(m Message, ws *websocket.Conn) {
	s.Keyboard <- m.Content[0]
	if len(s.Computer.Irq) < 1 {
		s.Computer.Irq <- true
	}
}

func (s *Server) Launch(logOut io.Writer) {
	globalServer = *s
	globalServer.Log = log.New(logOut, "Web Server: ", 0)

	asset_path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/lseelenbinder/armsim/web/assets/")
	globalServer.Log.Println(asset_path)

	http.Handle("/", http.FileServer(http.Dir(asset_path)))
	http.Handle("/ws", websocket.Handler(wsHandler))

	if err := http.ListenAndServe("localhost:4567", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func wsHandler(ws *websocket.Conn) {
	globalServer.Serve(ws)
}
