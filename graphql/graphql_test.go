package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var schema = `{
	"data": {
		"__schema": {
			"types": [
				{
          "name": "SearchRequest",
          "fields": null
        },
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

func TestProbe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			t.Errorf("Expected to request '/query', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(schema))
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

	if result["status"] != "ok" {
		t.Fatalf(`Expected to get {"status":"ok"}, got: %s`, result)
	}
}

func TestLookForFieldsThatSeemSensitive(t *testing.T) {
	var jsonSchema map[string]interface{}
	err := json.Unmarshal([]byte(schema), &jsonSchema)
	if err != nil {
		t.Fatalf("Test setup horked: %s", err)
	}

	fields := lookForFieldsThatSeemSensitive(jsonSchema)

	if len(fields) != 1 {
		t.Fatalf(`Expected to find 1 sensitive fields, got: %d`, len(fields))
	}
}
