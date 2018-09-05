package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
			"url":      "status",
			"response": "ok",
		},
		{
			"url":      "not-found",
			"response": "404 page not found\n",
		},
	}
}

func testRig(f func()) {
	server := newServer(port)
	server.listenAndServe()
	defer server.shutdown()
	f()
}

var host = "http://localhost:3333"

func testHttpRequest(verb string, resource string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	return client.Do(r)
}

func TestJsonRequest(t *testing.T) {
	testRig(func() {
		response, err := testHttpRequest("POST", "/name", `{"name": "Hodor"}`)
		if err != nil {
			t.Fatalf("Request failer %v", err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Error reading response body %v", err)
		}
		var r Response
		err = json.Unmarshal(body, &r)
		if err != nil {
			t.Fatalf("Error transforming response body to JSON %v. Body: %v", err, string(body))
		}
		if r.Msg != "Hodor" {
			t.Fatalf("The end-point do not respond correctly %v", r.Msg)
		}
	})
}
