package warcraft

import (
    "./io"
    "fmt"
)

type MapInfo struct {
    Dummy [0x14]byte
    Path string
    HostName string
}

type GameInfoPacket struct {
    ClientVersion ClientVersion
    Id uint32
    EntryKey uint32
    Name string
    MapInfo MapInfo `encode:"nested"`
    Slots uint32
    gameType [4]byte
    CurrentPlayers uint32
    PlayerSlots uint32
    UpTime uint32
    Port uint16
}

var GameInfoPacketType = byte(0x30)

func (gameInfoPacket *GameInfoPacket) PacketType() byte {
    return GameInfoPacketType
}

func init() {
    io.RegisterPacketType(GameInfoPacketType, func() io.Packet {
        return new(GameInfoPacket)
    })
}

func (gameInfo *GameInfoPacket) Bytes() []byte {
    return io.PacketBytes(gameInfo)
}

func ParseGameInfoPacket(data []byte) (gameInfo GameInfoPacket, err error) {
    err = io.ParsePacket(&gameInfo, data)
    if err != nil {
        err = fmt.Errorf("Unable to parse game info packet: %v", err)
    }
    return
}

func (game GameInfoPacket) String() string {
    return fmt.Sprintf("Name: %q, Map: %q", game.Name, game.MapInfo.Path)
}
