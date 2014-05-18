package warcraft

import (
    "./io"
    "fmt"
)

type CancelPacket struct {
    GameId uint32
}

var CancelPacketType = byte(0x33)

func (cancelPacket *CancelPacket) PacketType() byte {
    return CancelPacketType
}

func init() {
    io.RegisterPacketType(CancelPacketType, func() io.Packet {
        return new(CancelPacket)
    })
}

func (cancel *CancelPacket) Bytes() []byte {
    return io.PacketBytes(cancel)
}

func ParseCancelPacket(data []byte) (cancel CancelPacket, err error) {
    err = io.ParsePacket(&cancel, data)
    if err != nil {
        err = fmt.Errorf("Unable to parse cancel packet: %v", err)
    }
    return
}

func NewCancelPacket(gameId uint32) CancelPacket {
    return CancelPacket{
        GameId: gameId,
    }
}
