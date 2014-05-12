package warcraft

import "testing"

func TestBrowsePacket(t *testing.T) {
    browse := NewBrowsePacket(NewClientVersion(RawExpansion, 19))
    b2, err := ParseBrowsePacket(browse.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if b2 != browse {
        t.Fatalf("b2 != browse")
    }
}
