package warcraft

import (
    "bytes"
    "testing"
)

func TestPacket(t *testing.T) {
    var buffer bytes.Buffer

    expectedPacketType := byte(10)
    expectedData := []byte { 1, 2, 3, 4, 5 }
    expectedRaw := []byte { PacketHeader, 10, 9, 0, 1, 2, 3, 4, 5 }

    WritePacket(&buffer, expectedPacketType, expectedData)
    if !bytes.Equal(buffer.Bytes(), expectedRaw) {
        t.Errorf("Failed")
    }

    packetType, data, err := ReadPacket(&buffer)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if packetType != expectedPacketType {
        t.Errorf("Failed: %d != %d", packetType, expectedPacketType)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }

    raw := PacketBytes(expectedPacketType, expectedData)
    if !bytes.Equal(raw, expectedRaw) {
        t.Errorf("Failed: %v != %v", raw, expectedRaw)
    }

    packetType, data, err = ParsePacket(expectedRaw)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if packetType != expectedPacketType {
        t.Errorf("Failed: %d != %d", packetType, expectedPacketType)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }
}
