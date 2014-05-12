package warcraft

import "testing"

func TestBrowse(t *testing.T) {
    browse := NewBrowse(NewClientVersion(RawExpansion, 19))
    b2, err := ParseBrowse(browse.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if b2 != browse {
        t.Fatalf("b2 != browse")
    }
}
