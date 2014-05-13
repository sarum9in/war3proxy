package warcraft

import "testing"

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
