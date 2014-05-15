package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type CancelPacket struct {
    GameId uint32
}

var CancelPacketHeader = [...]byte { 0xf7, 0x33, 0x08, 0x00 }
const CancelPacketSize = 8

func NewCancelPacket(gameId uint32) CancelPacket {
    return CancelPacket{
        GameId: gameId,
    }
}

func (CancelPacket *CancelPacket) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(CancelPacketHeader[:])
    binary.Write(&buffer, binary.LittleEndian, CancelPacket.GameId)

    result := buffer.Bytes()
    if len(result) != CancelPacketSize {
        panic(fmt.Errorf("len(result) != CancelPacketSize"))
    }
    return result
}

func ParseCancelPacket(data []byte) (CancelPacket CancelPacket, err error) {
    if len(data) != CancelPacketSize {
        err = fmt.Errorf("len(data) != CancelPacketSize")
        return
    }

    if !bytes.HasPrefix(data, CancelPacketHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    buffer := bytes.NewBuffer(data[4:])

    err = binary.Read(buffer, binary.LittleEndian, &CancelPacket.GameId)
    if err != nil {
        err = fmt.Errorf("Unable to parse CancelPacket.GameId: %v", err)
        return
    }

    return
}
