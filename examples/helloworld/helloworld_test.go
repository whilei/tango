package main

import (
	"testing"

	"github.com/unrolled/tango"
)

func TestResponse(t *testing.T) {
	client := tango.NewTestClient(t)

	resp := client.Get("/")

	if resp.StatusCode != 200 {
		t.Errorf("resp.StatusCode(%v) != 200", resp.StatusCode)
		t.FailNow()
	}

	if resp.Content != "Hello, world." {
		t.Errorf("resp.Content('%v') != 'Hello, world.'", resp.Content)
		t.FailNow()
	}
}
