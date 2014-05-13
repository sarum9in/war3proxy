package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type BrowsePacket struct {
    ClientVersion ClientVersion
}

var BrowsePacketHeader = [...]byte{ 0xf7, 0x2f, 0x10, 0x00 }
const BrowsePacketSize = 16

func NewBrowsePacket(clientVersion ClientVersion) BrowsePacket {
    return BrowsePacket{ ClientVersion: clientVersion }
}

func (browse *BrowsePacket) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(BrowsePacketHeader[:])
    buffer.Write(browse.ClientVersion.Expansion[:])
    binary.Write(&buffer, binary.LittleEndian, browse.ClientVersion.Version)
    dummy := [...]byte { 0x00, 0x00, 0x00, 0x00 }
    buffer.Write(dummy[:])

    result := buffer.Bytes()
    if len(result) != BrowsePacketSize {
        panic(fmt.Errorf("len(result) != BrowsePacketSize"))
    }
    return result
}

func ParseBrowsePacket(data []byte) (browse BrowsePacket, err error) {
    if len(data) != BrowsePacketSize {
        err = fmt.Errorf("len(data) != BrowsePacketSize")
        return
    }

    if !bytes.HasPrefix(data, BrowsePacketHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    copy(browse.ClientVersion.Expansion[:], data[4:8])

    err = binary.Read(bytes.NewBuffer(data[8:]), binary.LittleEndian, &browse.ClientVersion.Version)
    if err != nil {
        err = fmt.Errorf("Unable to parse ClientVersion.Expansion: %v", err)
        return
    }

    // ignore unknown 4-byte dummy

    return
}
