package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type AnnouncePacket struct {
    GameId uint32
    Players uint32
    Slots uint32
}

var AnnouncePacketHeader = [...]byte { 0xf7, 0x32, 0x10, 0x00 }
const AnnouncePacketSize = 16

func NewAnnouncePacket(gameId uint32, players uint32, slots uint32) AnnouncePacket {
    return AnnouncePacket{
        GameId: gameId,
        Players: players,
        Slots: slots,
    }
}

func (announce *AnnouncePacket) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(AnnouncePacketHeader[:])
    binary.Write(&buffer, binary.LittleEndian, announce.GameId)
    binary.Write(&buffer, binary.LittleEndian, announce.Players)
    binary.Write(&buffer, binary.LittleEndian, announce.Slots)

    result := buffer.Bytes()
    if len(result) != AnnouncePacketSize {
        panic(fmt.Errorf("len(result) != AnnouncePacketSize"))
    }
    return result
}

func ParseAnnouncePacket(data []byte) (announce AnnouncePacket, err error) {
    if len(data) != AnnouncePacketSize {
        err = fmt.Errorf("len(data) != AnnouncePacketSize")
        return
    }

    if !bytes.HasPrefix(data, AnnouncePacketHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    buffer := bytes.NewBuffer(data[4:])

    err = binary.Read(buffer, binary.LittleEndian, &announce.GameId)
    if err != nil {
        err = fmt.Errorf("Unable to parse AnnouncePacket.GameId: %v", err)
        return
    }

    err = binary.Read(buffer, binary.LittleEndian, &announce.Players)
    if err != nil {
        err = fmt.Errorf("Unable to parse AnnouncePacket.Players: %v", err)
        return
    }

    err = binary.Read(buffer, binary.LittleEndian, &announce.Slots)
    if err != nil {
        err = fmt.Errorf("Unable to parse AnnouncePacket.Slots: %v", err)
        return
    }

    return
}
