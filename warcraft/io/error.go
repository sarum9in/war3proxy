package io

import "fmt"

type UnregisteredPacketTypeError struct {
    PacketType byte
}

func (e *UnregisteredPacketTypeError) Error() string {
    return fmt.Sprintf(
        "Unable to create packet with unregistered type = %x",
        e.PacketType,
    )
}

type UnexpectedPacketTypeError struct {
    PacketType         byte
    ExpectedPacketType byte
}

func (e *UnexpectedPacketTypeError) Error() string {
    return fmt.Sprintf("Unexpected packet type %x, expected %x",
        e.PacketType, e.ExpectedPacketType)
}

type UnexpectedBigPacket struct {
    Size int
}

func (e *UnexpectedBigPacket) Error() string {
    return fmt.Sprintf("Unexpected big packet size = ", e.Size)
}

type InvalidPacketHeaderError struct {
    Header         byte
    ExpectedHeader byte
}

func (e *InvalidPacketHeaderError) Error() string {
    return fmt.Sprintf("Invalid header %x, expected %x",
        e.Header, e.ExpectedHeader)
}

type InvalidPacketDataSizeError struct {
    ActualSize   int
    DeclaredSize int
}

func (e *InvalidPacketDataSizeError) Error() string {
    return fmt.Sprintf("Invalid packet data size %d, declared %d",
        e.ActualSize, e.DeclaredSize)
}

type UnexpectedNullByteError struct {
    Data []byte
    Pos  int
}

func (e *UnexpectedNullByteError) Error() string {
    return fmt.Sprintf("Unexpected null byte at position = %d", e.Pos)
}
