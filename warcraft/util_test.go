package warcraft

import (
    "bytes"
    "testing"
)

func TestInteger(t *testing.T) {
    var buffer bytes.Buffer

    var integer uint32 = 10
    WriteInteger(&buffer, integer)

    var integer2 uint32
    ReadInteger(&buffer, &integer2)
    if integer2 != integer {
        t.Errorf("Failed: %d != %d", integer2, integer)
    }
}

func TestNullTerminated(t *testing.T) {
    var buffer bytes.Buffer

    WriteNullTerminatedBytes(&buffer, []byte {1, 2, 3, 4, 5})
    if !bytes.Equal(buffer.Bytes(), []byte {1, 2, 3, 4, 5, 0}) {
        t.Errorf("Failed write")
    }

    if !bytes.Equal(ReadNullTerminatedBytes(&buffer), []byte {1, 2, 3, 4, 5}) {
        t.Errorf("Failed read")
    }
}

func TestEncoding(t *testing.T) {
    tests := [...]string {
        "",
        "0",
        "01",
        "012",
        "0123",
        "01234",
        "012345",
        "0123456",
        "01234567",
        "012345678",
        "0123456789",
        "01234567890",
        "012345678901",
        "0123456789012",
        "01234567890123",
        "012345678901234",
        "0123456789012345",
        "01234567890123456",
        "012345678901234567",
        "0123456789012345678",
        "01234567890123456789",
        "hello",
        "hello, world",
    }
    for _, s := range tests {
        bs := []byte(s)
        ebs := EncodeBytes(bs)
        debs := DecodeBytes(ebs)
        sdebs := string(debs)
        if sdebs != s {
            t.Errorf("Failed: %q != %q", sdebs, s)
        }
    }
}
