package warcraft

import (
	"bytes"
	"errors"
	"encoding/binary"
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

    cryptstart := 0x14 + nameLength + 1
    decrypted := DecryptString(data[cryptstart:])
    game.Map, _ = ParseString(decrypted[0xd:])

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

func ParseInteger(data []byte, integer interface{}) {
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, integer)
    if err != nil {
        panic(err)
    }
}

func ParseString(data []byte) (str string, rawSize int) {
    raw := ParseRawString(data)
    str = string(raw) // UTF8
    rawSize = len(raw) + 1
	return
}

func ParseRawString(data []byte) []byte {
	var output bytes.Buffer
	for _, c := range data {
		if c == 0 {
			break
		} else {
			output.WriteByte(byte(c))
		}
	}
	return output.Bytes()
}

func DecryptString(data []byte) []byte {
	var output bytes.Buffer
    var mask byte = 0
    for pos, c := range data {
        if pos % 8 == 0 {
            mask = c
        } else {
            if mask & (0x1 << uint(pos % 8)) == 0 {
                output.WriteByte(c - 1)
            } else {
                output.WriteByte(c)
            }
        }
    }
    return output.Bytes()
}
