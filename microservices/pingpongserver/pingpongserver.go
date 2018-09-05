package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

type Response struct {
	Msg string `json:"msg,omitempty"`
}

type Request struct {
	Name string `json:"name"`
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
	nameHandle := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading request body", err)
			return
		}

		var request Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Println("decode the incoming json fail", err)
			return
		}

		log.Println("Request: ", request)
		response := &Response{Msg: request.Name}
		b, err := json.Marshal(response)
		if err != nil {
			log.Println("json marshal error", err)
		}
		stringResponse := string(b)
		log.Println("Response: ", stringResponse)
		io.WriteString(w, stringResponse)
	}

	status := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/name", nameHandle)
	mux.HandleFunc("/status", status)

	httpServer := &http.Server{Addr: ":" + port, Handler: mux}
	return &server{httpServer: httpServer}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "3333", "./pingpongserver -port 3333")
	flag.Parse()

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
	fmt.Println("POST: :PORT/name {name: \"Andrea\"}")
	<-done
	fmt.Println("Exiting...")
}
