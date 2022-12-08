package main

import (
	"testing"
)

func TestLogin(t *testing.T) {
	var localhost string = "127.0.0.1"
	redisInstance, err := GetKeys(localhost)

	if err != nil && redisInstance.Name != "" {
		t.Fatalf(`Login(localhost) result was %s, want empty string`, redisInstance.Name)
	}

	// need to setup context to test this
	if redisInstance.Name == "redis" && err != nil {
		t.Fatalf(`Login("localhost") error was %v, want nil`, err)
	}
}
