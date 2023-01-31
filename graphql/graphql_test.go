package main

import (
	"testing"
)

func TestProbe(t *testing.T) {
	err := Probe("http://localhost:8080/query")

	if err != nil {
		t.Fatalf(`Probing horked`)
	}
}
