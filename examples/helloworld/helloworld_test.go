package main

import (
    "github.com/cojac/assert"
    "github.com/cojac/tango"
    "testing"
)

func TestResponse(t *testing.T) {
    client := tango.NewTestClient(t)

    resp := client.Get("/")

    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "Hello, world.", resp.Content)
}
