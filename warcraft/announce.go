package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type Announce struct {
    GameId uint32
    Players uint32
    Slots uint32
}

var AnnounceHeader = [...]byte { 0xf7, 0x32, 0x10, 0x00 }
const AnnounceSize = 16

func NewAnnounce(gameId uint32, players uint32, slots uint32) Announce {
    return Announce{
        GameId: gameId,
        Players: players,
        Slots: slots,
    }
}

func (announce *Announce) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(AnnounceHeader[:])
    binary.Write(&buffer, binary.LittleEndian, announce.GameId)
    binary.Write(&buffer, binary.LittleEndian, announce.Players)
    binary.Write(&buffer, binary.LittleEndian, announce.Slots)

    result := buffer.Bytes()
    if len(result) != AnnounceSize {
        panic(fmt.Errorf("len(result) != AnnounceSize"))
    }
    return result
}

func ParseAnnounce(data []byte) (announce Announce, err error) {
    if len(data) != AnnounceSize {
        err = fmt.Errorf("len(data) != AnnounceSize")
        return
    }

    if !bytes.HasPrefix(data, AnnounceHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    buffer := bytes.NewBuffer(data[4:])

    err = binary.Read(buffer, binary.LittleEndian, &announce.GameId)
    if err != nil {
        err = fmt.Errorf("Unable to parse Announce.GameId: %v", err)
        return
    }

    err = binary.Read(buffer, binary.LittleEndian, &announce.Players)
    if err != nil {
        err = fmt.Errorf("Unable to parse Announce.Players: %v", err)
        return
    }

    err = binary.Read(buffer, binary.LittleEndian, &announce.Slots)
    if err != nil {
        err = fmt.Errorf("Unable to parse Announce.Slots: %v", err)
        return
    }

    return
}
