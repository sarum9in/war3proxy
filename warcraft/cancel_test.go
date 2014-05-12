package warcraft

import "testing"

func TestCancelPacket(t *testing.T) {
    CancelPacket := NewCancelPacket(11)
    a2, err := ParseCancelPacket(CancelPacket.Bytes())
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    if a2 != CancelPacket {
        t.Fatalf("a2 != CancelPacket")
    }
}
