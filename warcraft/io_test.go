package warcraft

import (
    "bytes"
    "testing"
)

func TestPacket(t *testing.T) {
    var buffer bytes.Buffer

    expectedPacket := byte(10)
    expectedData := []byte { 1, 2, 3, 4, 5 }
    expectedRaw := []byte { PacketHeader, 10, 9, 0, 1, 2, 3, 4, 5 }

    WritePacket(&buffer, expectedPacket, expectedData)
    if !bytes.Equal(buffer.Bytes(), expectedRaw) {
        t.Errorf("Failed")
    }

    packet, data, err := ReadPacket(&buffer)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if packet != expectedPacket {
        t.Errorf("Failed: %d != %d", packet, expectedPacket)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }

    raw := PacketBytes(expectedPacket, expectedData)
    if !bytes.Equal(raw, expectedRaw) {
        t.Errorf("Failed: %v != %v", raw, expectedRaw)
    }

    packet, data, err = ParsePacket(expectedRaw)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if packet != expectedPacket {
        t.Errorf("Failed: %d != %d", packet, expectedPacket)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }
}
