package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type Cancel struct {
    GameId uint32
}

var CancelHeader = [...]byte { 0xf7, 0x33, 0x08, 0x00 }
const CancelSize = 8

func NewCancel(gameId uint32) Cancel {
    return Cancel{
        GameId: gameId,
    }
}

func (Cancel *Cancel) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(CancelHeader[:])
    binary.Write(&buffer, binary.LittleEndian, Cancel.GameId)

    result := buffer.Bytes()
    if len(result) != CancelSize {
        panic(fmt.Errorf("len(result) != CancelSize"))
    }
    return result
}

func ParseCancel(data []byte) (Cancel Cancel, err error) {
    if len(data) != CancelSize {
        err = fmt.Errorf("len(data) != CancelSize")
        return
    }

    if !bytes.HasPrefix(data, CancelHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    buffer := bytes.NewBuffer(data[4:])

    err = binary.Read(buffer, binary.LittleEndian, &Cancel.GameId)
    if err != nil {
        err = fmt.Errorf("Unable to parse Cancel.GameId: %v", err)
        return
    }

    return
}

