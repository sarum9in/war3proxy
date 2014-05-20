package io

import (
    "bytes"
    "io"
    "testing"
)

func TestNested(t *testing.T) {
    tests := [...]string{
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
        ebs := EncodeNestedBytes(bs)
        if ebs[len(ebs)-1] != 0 {
            t.Errorf("Failed: last byte %d != %d", ebs[len(ebs)-1], 0)
        }
        debs := DecodeNestedBytes(ebs)
        sdebs := string(debs)
        if sdebs != s {
            t.Errorf("Failed: %q != %q", sdebs, s)
        }
    }

    expectedData := []byte{
        1, 2, 3, 4, 0, 5, 6,
        7, 8, 0, 9, 10, 11, 12,
        0, 13, 14,
    }
    expectedEncoded := []byte{
        75, 1, 3, 3, 5, 1, 5, 7,
        83, 7, 9, 1, 9, 11, 11, 13,
        5, 1, 13, 15,
        0,
    }
    var buffer bytes.Buffer
    writer := NewNestedWriter(&buffer)
    _, err := writer.Write(expectedData)
    writer.Close()
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if !bytes.Equal(buffer.Bytes(), expectedEncoded) {
        t.Errorf("Failed: %v != %v", buffer.Bytes(), expectedEncoded)
    }

    data := make([]byte, 4096)
    reader := NewNestedReader(&buffer)
    n, err := reader.Read(data)
    data = data[:n]
    if err != nil && err != io.EOF {
        t.Errorf("Failed: %v", err)
    }
    if !bytes.Equal(data, expectedData) {
        t.Errorf("Failed: %v != %v", data, expectedData)
    }
}
