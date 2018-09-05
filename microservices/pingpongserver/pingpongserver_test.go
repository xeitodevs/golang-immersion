package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServerPass(t *testing.T) {
	testRig(func() {
		for _, data := range getDataProvider() {
			res, err := http.Get("http://localhost:3333/" + data["url"])
			if err != nil {
				t.Fatal(err)
			}

			msg, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			if string(msg) != data["response"] {
				t.Fatal("Unexpected message", string(msg))
			}
		}
	})
}

func getDataProvider() [2]map[string]string {
	return [2]map[string]string{
		{
			"url":      "pong",
			"response": "Ping",
		},
		{
			"url":      "ping",
			"response": "Pong",
		},
	}
}

func testRig(f func()) {
	server := newServer(port)
	server.listenAndServe()
	defer server.shutdown()
	f()
}
