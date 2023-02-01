package main

import (
	"testing"
)

func TestLogin(t *testing.T) {
	var localhost string = "127.0.0.1"
	clusterDetails, err := Login(localhost)

	if err != nil && clusterDetails.Name != "" {
		t.Fatalf(`Login(localhost) result was %s, want empty string`, clusterDetails.Name)
	}

	// need to setup context to test this
	if clusterDetails.Name == "elasticsearch" && err != nil {
		t.Fatalf(`Login("localhost") error was %v, want nil`, err)
	}
}
