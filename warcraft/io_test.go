package warcraft

import (
    "bytes"
    "testing"
)

func TestPacket(t *testing.T) {
    var buffer bytes.Buffer

    expectedData := []byte { 1, 2, 3, 4, 5 }

    WritePacket(&buffer, 10, expectedData)
    if !bytes.Equal(buffer.Bytes(), []byte { PacketHeader, 10, 9, 0, 1, 2, 3, 4, 5 }) {
        t.Errorf("Failed")
    }

    packet, data, err := ReadPacket(&buffer)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if packet != 10 {
        t.Errorf("Failed: %d != %d", packet, 10)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }
}
