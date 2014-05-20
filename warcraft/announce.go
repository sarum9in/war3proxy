package warcraft

import "./io"

type AnnouncePacket struct {
    GameId  uint32
    Players uint32
    Slots   uint32
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

func (announce *AnnouncePacket) Parse(data []byte) (err error) {
    err = io.ParsePacket(announce, data)
    if err != nil {
        err = &ParseError{
            Name: "announce packet",
            Err:  err,
        }
    }
    return
}

func ParseAnnouncePacket(data []byte) (announce AnnouncePacket, err error) {
    err = announce.Parse(data)
    return
}

func NewAnnouncePacket(gameId uint32, players uint32, slots uint32) AnnouncePacket {
    return AnnouncePacket{
        GameId:  gameId,
        Players: players,
        Slots:   slots,
    }
}
