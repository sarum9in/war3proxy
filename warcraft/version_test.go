package warcraft

import "testing"

func TestClientVersion(t *testing.T) {
    a := NewClientVersion(RawExpansion, 20)
    b := ClientVersion{
        Expansion: RawExpansion,
        Version: 20,
    }

    if a != b {
        t.Error()
    }
}
