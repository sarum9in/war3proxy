package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
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

func ParsePacket(data []byte) (packet byte, packetData []byte, err error) {
    packet, packetData, err = ReadPacket(bytes.NewReader(data))
    if err != nil {
        return
    }

    expectedPacketDataSize := len(data) - 2 - 2
    if len(packetData) != expectedPacketDataSize {
        err = fmt.Errorf("Invalid packet's data size: %d != %d", len(packetData), expectedPacketDataSize)
        return
    }

    return
}

func PacketBytes(packet byte, data []byte) []byte {
    var buffer bytes.Buffer

    err := WritePacket(&buffer, packet, data)
    if err != nil {
        panic(err)
    }

    return buffer.Bytes()
}

func ReadPacket(reader Reader) (packet byte, data []byte, err error) {
    var header [2]byte
    _, err = reader.Read(header[:])
    if err != nil {
        return
    }

    if header[0] != PacketHeader {
        err = fmt.Errorf("Invalid header")
        return
    }

    packet = header[1]

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

func WritePacket(writer Writer, packet byte, data []byte) (err error) {
    err = writer.WriteByte(PacketHeader)
    if err != nil {
        return
    }

    err = writer.WriteByte(packet)
    if err != nil {
        return
    }

    realSize := 2 + 2 + len(data) // header + size + data
    var size uint16 = uint16(realSize)
    if int(size) != realSize {
        err = fmt.Errorf("Too big packet: size = %d", realSize)
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
