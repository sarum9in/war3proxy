package warcraft

import (
    "bytes"
    "testing"
)

func TestReflect(t *testing.T) {
    type MyStruct struct {
        Ten uint32
        Dummy [5]byte
        Hello string
        World string
    }

    expectedPacketData := []byte { 10, 0, 0, 0, 1, 2, 3, 4, 5, 104, 101, 108, 108, 111, 0, 119, 111, 114, 108, 100, 0 }
    expectedDummy := []byte { 1, 2, 3, 4, 5 }

    reader := bytes.NewReader(expectedPacketData)

    data := new(MyStruct)
    err := ReflectRead(reader, data)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if data.Ten != 10 {
        t.Errorf("Failed: %d != %d", data.Ten, 10)
    }
    if !bytes.Equal(data.Dummy[:], expectedDummy) {
        t.Errorf("Failed: %v != %v", data.Dummy, expectedDummy)
    }
    if data.Hello != "hello" {
        t.Errorf("Failed: %q != %q", data.Hello, "hello")
    }
    if data.World != "world" {
        t.Errorf("Failed: %q != %q", data.World, "world")
    }

    var writer bytes.Buffer

    err = ReflectWrite(&writer, data)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if !bytes.Equal(writer.Bytes(), expectedPacketData) {
        t.Errorf("Failed: %v != %v", writer.Bytes(), expectedPacketData)
    }
}

func TestReflectWrite(t *testing.T) {
}
