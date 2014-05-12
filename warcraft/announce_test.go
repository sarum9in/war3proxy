package warcraft

import "testing"

func TestAnnounce(t *testing.T) {
    announce := NewAnnounce(12, 20, 45)
    a2, err := ParseAnnounce(announce.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if a2 != announce {
        t.Fatalf("a2 != announce")
    }
}
