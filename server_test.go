package tango

import (
    "github.com/cojac/assert"
    "testing"
)

func TestServerVersion(t *testing.T) {
    if len(Version) < 5 {
        t.Errorf("Version should be at least 5 chars.")
    }
}

func TestServerVersionMap(t *testing.T) {
    oldVer := Version
    Version = "2.5.11"

    vmap := VersionMap()
    assert.Equal(t, 2, vmap[0])
    assert.Equal(t, 5, vmap[1])
    assert.Equal(t, 11, vmap[2])

    Version = oldVer
}
