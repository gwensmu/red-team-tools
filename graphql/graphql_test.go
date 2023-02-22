package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			t.Errorf("Expected to request '/query', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	err := Probe(server.URL + "/query")

	if err != nil {
		t.Fatalf(`Probing horked`)
	}
}

func TestSendQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			t.Errorf("Expected to request '/query', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	result, err := sendQuery(server.URL+"/query", "hello")

	if err != nil {
		t.Fatalf(`Probing horked`)
	}

	if result != `{"status":"ok"}` {
		t.Fatalf(`Expected to get {"status":"ok"}, got: %s`, result)
	}
}

func TestLookForFieldsThatSeemSensitive(t *testing.T) {
	schema := `{
		"data": {
			"__schema": {
				"types": [
					{
						"name": "User",
						"fields": [
							{
								"name": "id",
								"args": []
							},
							{
								"name": "ssn",
								"args": []
							}
						]
					}
				]
			}
		}
	}`

	fields := lookForFieldsThatSeemSensitive(schema)

	if len(fields) != 1 {
		t.Fatalf(`Expected to find 1 sensitive fields, got: %d`, len(fields))
	}
}
