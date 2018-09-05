package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	// a channel to receive unix signals
	sigs := make(chan os.Signal, 1)
	// a channel to receive a stop confirmation on interrupt
	done := make(chan bool, 1)
	// signal.Notify is a method to create a channel which receives
	// SIGINT, SIGTERM unix signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	moveAlong := func() {
		fmt.Println("Not the droid you lookin for...")
	}

	server := newServer(port)
	server.listenAndServe()
	defer moveAlong()

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		server.shutdown()
		done <- true
	}()
	// Ctrl-C sends a SIGINT signal to the program
	fmt.Println("Ctrl-C to interrupt...")
	<-done
	fmt.Println("Exiting...")
}
