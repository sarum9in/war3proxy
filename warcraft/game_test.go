package warcraft

import (
    "bytes"
    "./io"
    "testing"
)

func TestGameInfoPacket(t *testing.T) {
    expectedGameInfo := GameInfoPacket{
        ClientVersion: ClientVersion{
            Expansion: TftExpansion,
            Version: 22,
        },
        Id: 10,
        EntryKey: 152,
        Name: "game name",
        Dummy: [0x01]byte { 0x07 },
        MapInfo: MapInfo{
            Dummy: [0x0d]byte {
                0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
                0x09, 0x0a, 0x0b, 0x0c, 0x0d,
            },
            Path: "some/path",
            HostName: "host name",
        },
        Slots: 10,
        GameType: [0x4]byte { 0x01, 0x02, 0x03, 0x04 },
        CurrentPlayers: 4,
        PlayerSlots: 6,
        UpTime: 20,
        Port: 6112,
    }

    expectedGameInfoData := []byte {
        0xf7, 0x30, 0x5c, 0x00, // header
        0x50, 0x58, 0x33, 0x57, 0x16, 0x00, 0x00, 0x00, // ClientVersion
        0x0a, 0x00, 0x00, 0x00, // Id
        0x98, 0x00, 0x00, 0x00, // EntryKey
        0x67, 0x61, 0x6d, 0x65, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x0, // Name
        0x07, // Dummy

        // MapInfo
        0xab, 0x01, 0x03, 0x03, 0x05, 0x05, 0x07, 0x07,
        0xd5, 0x09, 0x09, 0x0b, 0x0b, 0x0d, 0x0d, 0x73,
        0x5f, 0x6f, 0x6d, 0x65, 0x2f, 0x71, 0x61, 0x75,
        0x31, 0x69, 0x01, 0x69, 0x6f, 0x73, 0x75, 0x21,
        0x1d, 0x6f, 0x61, 0x6d, 0x65, 0x01,
        0x00,

        0x0a, 0x00, 0x00, 0x00, // Slots
        0x01, 0x02, 0x03, 0x04, // GameType
        0x04, 0x00, 0x00, 0x00, // CurrentPlayers
        0x06, 0x00, 0x00, 0x00, // PlayerSlots
        0x14, 0x00, 0x00, 0x00, // UpTime
        0xe0, 0x17, // Port
    }

    expectedMapInfoData := []byte {
        // Dummy
        0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
        0x09, 0x0a, 0x0b, 0x0c, 0x0d,

        0x73, 0x6f, 0x6d, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x00, // Path
        0x68, 0x6f, 0x73, 0x74, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x00, // HostName
    }

    var err error

    var buffer bytes.Buffer

    err = io.ReflectWrite(&buffer, expectedGameInfo.MapInfo)
    if err != nil {
        t.Logf("Failed: %v", err)
    }
    if !bytes.Equal(buffer.Bytes(), expectedMapInfoData) {
        t.Errorf("Failed: %v != %v", buffer.Bytes(), expectedMapInfoData)
    }

    var mapInfo MapInfo
    err = io.ReflectRead(&buffer, &mapInfo)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if mapInfo != expectedGameInfo.MapInfo {
        t.Errorf("Failed: %v != %v", mapInfo, expectedGameInfo.MapInfo)
    }

    data := expectedGameInfo.Bytes()
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if !bytes.Equal(data, expectedGameInfoData) {
        t.Errorf("Failed: %#v != %v", data, expectedGameInfoData)
    }

    gameInfo, err := ParseGameInfoPacket(data)
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if gameInfo != expectedGameInfo {
        t.Errorf("Failed: %#v != %#v", gameInfo, expectedGameInfo)
    }
}
