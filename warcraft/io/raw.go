package io

import (
    "bytes"
    "encoding/binary"
    "io"
)

type Reader interface {
    io.Reader
    io.ByteReader
}

type Writer interface {
    io.Writer
    io.ByteWriter
}

var PacketHeader byte = 0xf7

func ParseRawPacket(data []byte) (packetType byte, packetData []byte, err error) {
    packetType, packetData, err = ReadRawPacket(bytes.NewReader(data))
    if err != nil {
        return
    }

    declaredPacketSize := len(packetData) + 2 + 2 // + header + size
    if len(data) != declaredPacketSize {
        err = &InvalidPacketDataSizeError{
            ActualSize:   len(data),
            DeclaredSize: declaredPacketSize,
        }
        return
    }

    return
}

func RawPacketBytes(packetType byte, data []byte) []byte {
    var buffer bytes.Buffer

    err := WriteRawPacket(&buffer, packetType, data)
    if err != nil {
        panic(err)
    }

    return buffer.Bytes()
}

func ReadRawPacket(reader Reader) (packetType byte, data []byte, err error) {
    var header [2]byte
    _, err = reader.Read(header[:])
    if err != nil {
        return
    }

    if header[0] != PacketHeader {
        err = &InvalidPacketHeaderError{
            Header:         header[0],
            ExpectedHeader: PacketHeader,
        }
        return
    }

    packetType = header[1]

    var size uint16
    ReadInteger(reader, &size)

    size -= 2 // header
    size -= 2 // size

    data = make([]byte, size)
    n, err := reader.Read(data)
    if err != nil {
        data = data[:n]
        return
    }

    return
}

func WriteRawPacket(writer Writer, packetType byte, data []byte) (err error) {
    err = writer.WriteByte(PacketHeader)
    if err != nil {
        return
    }

    err = writer.WriteByte(packetType)
    if err != nil {
        return
    }

    actualSize := 2 + 2 + len(data) // header + size + data
    var size uint16 = uint16(actualSize)
    if int(size) != actualSize {
        err = &UnexpectedBigPacket{
            Size: actualSize,
        }
        return
    }
    err = WriteInteger(writer, size)
    if err != nil {
        return
    }

    _, err = writer.Write(data)
    if err != nil {
        return
    }

    return
}

func ReadInteger(reader Reader, integer interface{}) error {
    return binary.Read(reader, binary.LittleEndian, integer)
}

func WriteInteger(writer Writer, integer interface{}) error {
    return binary.Write(writer, binary.LittleEndian, integer)
}

func ReadNullTerminatedBytes(reader Reader) (data []byte, err error) {
    var buffer bytes.Buffer

    for c, err := reader.ReadByte(); err == nil; c, err = reader.ReadByte() {
        if c == 0 {
            break
        } else {
            buffer.WriteByte(byte(c))
        }
    }
    if err != nil && err != io.EOF {
        return
    }

    data = buffer.Bytes()

    return
}

func WriteNullTerminatedBytes(writer Writer, data []byte) (err error) {
    for pos, c := range data {
        if c == 0 {
            panic(&UnexpectedNullByteError{
                Data: data,
                Pos:  pos,
            })
        }
    }

    _, err = writer.Write(data)
    if err != nil {
        return
    }

    err = writer.WriteByte(0)
    if err != nil {
        return
    }

    return
}

func ReadNullTerminatedString(reader Reader) (str string, err error) {
    raw, err := ReadNullTerminatedBytes(reader)
    if err != nil {
        return
    }

    str = string(raw) // UTF8

    return
}

func WriteNullTerminatedString(writer Writer, data string) (err error) {
    err = WriteNullTerminatedBytes(writer, []byte(data)) // UTF8
    if err != nil {
        return
    }

    return
}
