package net_helpers

import (
	"testing"
)

func TestGetGCEPrefixes(t *testing.T) {
	cidrs := GetGCEPrefixes("us-central1")

	if len(cidrs) < 1 {
		t.Fatalf(`GetGCEPrefixes("us-central1") was empty`)
	}
}
