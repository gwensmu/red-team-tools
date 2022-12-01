package main

import (
	"testing"
)

func TestHosts(t *testing.T) {
	hosts, err := Hosts("127.0.0.1/30")

	if hosts[0] != "127.0.0.1" || err != nil {
		t.Fatalf(`Hosts("127.0.0.1/30")[0] = %q, %v, want "127.0.0.1", nil`, hosts[0], err)
	}

	if hosts[1] != "127.0.0.2" || err != nil {
		t.Fatalf(`Hosts("127.0.0.1/30")[0] = %q, %v, want "127.0.0.2", nil`, hosts[0], err)
	}

	if len(hosts) != 2 {
		t.Fatalf(`Hosts("127.0.0.1/30") was length %d, want 2`, len(hosts))
	}
}

func TestDial(t *testing.T) {
	var localhost string = "127.0.0.1"
	result, err := Dial(localhost, 80)

	if err != nil && result != "" {
		t.Fatalf(`Dial(localhost) result was %s, want empty string`, result)
	}

	// need to setup context to test this
	if result == "ok" && err != nil {
		t.Fatalf(`Dial("localhost") error was %v, want nil`, err)
	}
}
func BenchmarkHosts(b *testing.B) {
}

func FuzzHosts(f *testing.F) {
	f.Skip("TODO")
}
