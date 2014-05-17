package warcraft

import (
    "errors"
    "fmt"
)

type GameInfo struct {
    Id uint32
    Name string
    Map string
    Slots uint32
    CurrentPlayers uint32
    PlayerSlots uint32
    Port uint16
}

func (game GameInfo) String() string {
    return fmt.Sprintf("Name: %q, Map: %q", game.Name, game.Map)
}

func ParseGameInfo(data []byte) (game GameInfo, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    if data[0] != 0xf7 || data[1] != 0x30 {
        err = errors.New("Not a game info")
        return
    }

    ParseInteger(data[0xc:], &game.Id)
    var nameLength int
    game.Name, nameLength = ParseString(data[0x14:])

    decoded := DecodeBytes(data[0x14 + nameLength + 1:])
    game.Map, _ = ParseString(decoded[0xd:])

    length := len(data)

    ParseInteger(data[length - 22:], &game.Slots)
    ParseInteger(data[length - 14:], &game.CurrentPlayers)
    ParseInteger(data[length - 10:], &game.PlayerSlots)
    ParseInteger(data[length - 2:], &game.Port)

    return
}

func ChangeServerPort(data []byte, port uint16) {
    length := len(data)
    data[length - 2] = byte(port & 0xff)
    data[length - 1] = byte((port >> 8) & 0xff)
}
