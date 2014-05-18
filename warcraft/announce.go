package warcraft

import (
    "./io"
    "fmt"
)

type AnnouncePacket struct {
    GameId uint32
    Players uint32
    Slots uint32
}

var AnnouncePacketType = byte(0x32)

func (announcePacket *AnnouncePacket) PacketType() byte {
    return AnnouncePacketType
}

func init() {
    io.RegisterPacketType(AnnouncePacketType, func() io.Packet {
        return new(AnnouncePacket)
    })
}

func (announce *AnnouncePacket) Bytes() []byte {
    return io.PacketBytes(announce)
}

func ParseAnnouncePacket(data []byte) (announce AnnouncePacket, err error) {
    err = io.ParsePacket(&announce, data)
    if err != nil {
        err = fmt.Errorf("Unable to parse announce packet: %v", err)
    }
    return
}

func NewAnnouncePacket(gameId uint32, players uint32, slots uint32) AnnouncePacket {
    return AnnouncePacket{
        GameId: gameId,
        Players: players,
        Slots: slots,
    }
}
