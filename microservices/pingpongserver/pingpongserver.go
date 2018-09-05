package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const port = "3333"

type server struct {
	httpServer *http.Server
	listener   net.Listener
}

func (s *server) listenAndServe() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return errors.New("Server error initialization")
	}
	s.listener = listener
	go s.httpServer.Serve(s.listener)
	fmt.Println("Server now listening")
	return nil
}

func (s *server) shutdown() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}
	fmt.Println("Shutting down server")
	return nil
}

func newServer(port string) *server {
	pongHandle := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Ping")
	}

	pingHandle := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Pong")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pong", pongHandle)
	mux.HandleFunc("/ping", pingHandle)

	httpServer := &http.Server{Addr: ":" + port, Handler: mux}
	return &server{httpServer: httpServer}
}

func main() {

	// Here comes the chan type, the Jackie Chan of concurrency(sorry). A "Channel"
	// is a typed pipe through which you can send and receive values across goroutines.
	// We are going to use channel as a block in the main goroutine.
	// This particular channel is of type struct{}. An empty struct occupies zero bytes
	// of storage and since we aren't going to actually send or receive any values it's used.
	// It's totally cool to have chan int, chan bool chan MyStruct etc.

	ch := make(chan struct{})

	moveAlong := func() {
		fmt.Println("Not the droid you lookin for...")
	}

	server := newServer(port)
	server.listenAndServe()
	defer server.shutdown()
	defer moveAlong()

	// This channel will wait to receive a value(and in our case it will wait for
	// eternity). While it's waiting further execution of the main goroutine will remain
	// blocked hence serving our purpose. To unblock this channel someone needs to
	// do : ch <- someVal or close(ch). More of this coming up.
	// Now you can do : go run simpleserver_2.go  and expect it to work.
	// Goto http://localhost:3333/another to check whether it did.
	// To exit : Ctrl-C works but that's not nice. We can be more graceful than that.
	<-ch
}
