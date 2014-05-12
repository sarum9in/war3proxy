package warcraft

import "testing"

func TestAnnouncePacket(t *testing.T) {
    announce := NewAnnouncePacket(12, 20, 45)
    a2, err := ParseAnnouncePacket(announce.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if a2 != announce {
        t.Fatalf("a2 != announce")
    }
}
