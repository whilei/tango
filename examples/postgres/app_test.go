package main

import (
	"strings"
	"testing"

	"github.com/unrolled/tango"
)

func TestResponse(t *testing.T) {
	client := tango.NewTestClient(t)

	// You can change settings like so;
	tango.Settings.Set("db_name", "postgres")
	tango.Settings.Set("db_user", "postgres")
	tango.Settings.Set("db_password", "")

	resp := client.Get("/")

	if resp.StatusCode != 200 {
		t.Errorf("resp.StatusCode(%v) != 200", resp.StatusCode)
		t.FailNow()
	}

	if !strings.Contains(resp.Content, "Postgres clock_timestamp is") {
		t.Errorf("resp.Content('%v') does not contain 'Postgres clock_timestamp is'", resp.Content)
		t.FailNow()
	}
}
