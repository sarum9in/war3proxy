package io

import (
    "bytes"
    "testing"
)

func TestReflect(t *testing.T) {
    type MyNestedStruct struct {
        Something string
        Five uint32
    }

    type MyStruct struct {
        Ten uint32
        Dummy [5]byte
        Nested MyNestedStruct `encoding:"nested"`
        Hello string
        World string
    }

    expectedPacketData := []byte {
        0x0a, 0x00, 0x00, 0x00, // Ten
        0x01, 0x02, 0x03, 0x04, 0x05, // Dummy

        // Nested
        0x9f, 0x73, 0x6f, 0x6d, 0x65, 0x75, 0x69, 0x69,
        0x15, 0x6f, 0x67, 0x01, 0x05, 0x01, 0x01, 0x01,
        0x00,

        0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0, // Hello
        0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0, // World
    }
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
    if data.Nested.Something != "something" {
        t.Errorf("Failed: %s != %s", data.Nested.Something, "something")
    }
    if data.Nested.Five != 5 {
        t.Errorf("Failed: %d != %d", data.Nested.Five, 5)
    }
    if data.Hello != "hello" {
        t.Errorf("Failed: %q != %q", data.Hello, "hello")
    }
    if data.World != "world" {
        t.Errorf("Failed: %q != %q", data.World, "world")
    }

    var writer bytes.Buffer

    err = ReflectWrite(&writer, *data)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if !bytes.Equal(writer.Bytes(), expectedPacketData) {
        t.Errorf("Failed: %v != %v", writer.Bytes(), expectedPacketData)
    }
}
