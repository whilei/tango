package tangoappengine

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

	if resp.Content != "Hello, appengine." {
		t.Errorf("resp.Content('%v') != 'Hello, appengine.'", resp.Content)
		t.FailNow()
	}
}
