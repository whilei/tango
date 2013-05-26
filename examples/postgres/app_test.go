package main

import (
    "github.com/cojac/assert"
    "github.com/cojac/tango"
    "testing"
)

func TestResponse(t *testing.T) {
    client := tango.NewTestClient(t)

    // You can change settings like so;
    tango.Settings.Set("db_name", "postgres")
    tango.Settings.Set("db_user", "postgres")
    tango.Settings.Set("db_password", "")

    resp := client.Get("/")

    assert.Equal(t, 200, resp.StatusCode)
    assert.Contains(t, "Postgres says the timestamp is:", resp.Content)
}
