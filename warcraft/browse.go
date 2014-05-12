package warcraft

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type Browse struct {
    ClientVersion ClientVersion
}

var BrowseHeader = [...]byte{ 0xf7, 0x2f, 0x10, 0x00 }
const BrowseSize = 12

func NewBrowse(clientVersion ClientVersion) Browse {
    return Browse{ ClientVersion: clientVersion }
}

func (browse *Browse) Bytes() []byte {
    var buffer bytes.Buffer

    buffer.Write(BrowseHeader[:])
    buffer.Write(browse.ClientVersion.Expansion[:])
    binary.Write(&buffer, binary.LittleEndian, browse.ClientVersion.Version)

    result := buffer.Bytes()
    if len(result) != BrowseSize {
        panic(fmt.Errorf("len(result) != BrowseSize"))
    }
    return result
}

func ParseBrowse(data []byte) (browse Browse, err error) {
    if len(data) != BrowseSize {
        err = fmt.Errorf("len(data) != BrowseSize")
        return
    }

    if !bytes.HasPrefix(data, BrowseHeader[:]) {
        err = fmt.Errorf("Invalid header")
        return
    }

    copy(browse.ClientVersion.Expansion[:], data[4:8])

    err = binary.Read(bytes.NewBuffer(data[8:]), binary.LittleEndian, &browse.ClientVersion.Version)
    if err != nil {
        err = fmt.Errorf("Unable to parse ClientVersion.Expansion: %v", err)
        return
    }

    return
}
