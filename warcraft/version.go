package warcraft

import "fmt"

type Expansion [4]byte

type ClientVersion struct {
	Expansion Expansion
	Version uint32
}

func (clientVersion ClientVersion) String() string {
    return fmt.Sprintf("Expansion: %q, Version: 1.%d",
                       string(clientVersion.Expansion[:]),
                       clientVersion.Version)
}

var RawExpansion Expansion = Expansion{0x33, 0x52, 0x41, 0x57} // 3RAW
var TftExpansion Expansion = Expansion{0x50, 0x58, 0x33, 0x57} // PX3W

func NewClientVersion(expansion Expansion, version uint32) ClientVersion {
    return ClientVersion{
        Expansion: expansion,
        Version: version,
    }
}
