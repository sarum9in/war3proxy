package warcraft

import "./io"

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

func (cancel *CancelPacket) Parse(data []byte) (err error) {
    err = io.ParsePacket(cancel, data)
    if err != nil {
        err = &ParseError{
            Name: "cancel packet",
            Err:  err,
        }
    }
    return
}

func ParseCancelPacket(data []byte) (cancel CancelPacket, err error) {
    err = cancel.Parse(data)
    return
}

func NewCancelPacket(gameId uint32) CancelPacket {
    return CancelPacket{
        GameId: gameId,
    }
}
