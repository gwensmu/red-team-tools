package main

import (
	"testing"
)

func TestHosts(t *testing.T) {
	hosts, err := Hosts("127.0.0.1/30")
	if err != nil {
		t.Error(err)
	}

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

func BenchmarkHosts(b *testing.B) {
}

func FuzzHosts(f *testing.F) {
	f.Skip("TODO")
}
