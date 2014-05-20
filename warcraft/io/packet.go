package io

import "bytes"

type Packet interface {
    PacketType() byte
}

type PacketInit func() Packet

var packetTypeMap = make(map[byte]PacketInit)

func RegisterPacketType(packetType byte, packetInit PacketInit) {
    packetTypeMap[packetType] = packetInit
}

func NewPacket(packetType byte) (packet Packet, err error) {
    packetInit, ok := packetTypeMap[packetType]
    if !ok {
        err = &UnregisteredPacketTypeError{
            PacketType: packetType,
        }
        return
    }

    packet = packetInit()

    return
}

func ParsePacket(packet Packet, data []byte) (err error) {
    reader := bytes.NewReader(data)

    err = PacketReadFrom(packet, reader)
    if err != nil {
        return err
    }

    return
}

func PacketBytes(packet Packet) []byte {
    var writer bytes.Buffer

    err := PacketWriteTo(packet, &writer)
    if err != nil {
        panic(err)
    }

    return writer.Bytes()
}

func PacketReadFrom(packet Packet, reader Reader) (err error) {
    packetType, data, err := ReadRawPacket(reader)
    if err != nil {
        return
    }
    if packetType != packet.PacketType() {
        err = &UnexpectedPacketTypeError{
            PacketType:         packetType,
            ExpectedPacketType: packet.PacketType(),
        }
        return
    }

    dataReader := bytes.NewReader(data)

    err = ReflectRead(dataReader, packet)
    if err != nil {
        return
    }

    return
}

func PacketWriteTo(packet Packet, writer Writer) (err error) {
    var dataWriter bytes.Buffer

    err = ReflectWrite(&dataWriter, packet)
    if err != nil {
        return
    }

    err = WriteRawPacket(writer, packet.PacketType(), dataWriter.Bytes())
    if err != nil {
        return
    }

    return
}

func ReflectReadPacket(reader Reader) (packet Packet, err error) {
    packetType, data, err := ReadRawPacket(reader)
    if err != nil {
        return
    }

    dataReader := bytes.NewReader(data)

    packet, err = NewPacket(packetType)
    if err != nil {
        return
    }

    err = ReflectRead(dataReader, packet)
    if err != nil {
        return
    }

    return
}
