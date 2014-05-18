package warcraft

import (
    "./io"
    "fmt"
)

type BrowsePacket struct {
    ClientVersion ClientVersion
    Dummy [4]byte
}

var BrowsePacketType = byte(0x2f)

func (browsePacket *BrowsePacket) PacketType() byte {
    return BrowsePacketType
}

func init() {
    io.RegisterPacketType(BrowsePacketType, func() io.Packet {
        return new(BrowsePacket)
    })
}

func (browse *BrowsePacket) Bytes() []byte {
    return io.PacketBytes(browse)
}

func ParseBrowsePacket(data []byte) (browse BrowsePacket, err error) {
    err = io.ParsePacket(&browse, data)
    if err != nil {
        err = fmt.Errorf("Unable to parse browse packet: %v", err)
    }
    return
}

func NewBrowsePacket(clientVersion ClientVersion) BrowsePacket {
    return BrowsePacket{ ClientVersion: clientVersion }
}
