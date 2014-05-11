package proxy

import (
    "bytes"
    "log"
    "net"
    "time"
)

func SendBrowse(conn *net.UDPConn, remote *net.UDPAddr, clientVersion ClientVersion) {
    data := []byte {
        0xf7, 0x2f, 0x10, 0x00,
        clientVersion.Expansion[0],
        clientVersion.Expansion[1],
        clientVersion.Expansion[2],
        clientVersion.Expansion[3],
        clientVersion.Version,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    }

    log.Println("Sending browse for client version:", clientVersion)
    _, err := conn.WriteTo(data, remote)
    if err != nil {
        log.Println("Unable to send browse:", err)
    }
}

func SendCancel(conn *net.UDPConn, game *GameInfo) {
    // FIXME int32 like a byte
    data := []byte {
        0xf7, 0x33, 0x08, 0x00, byte(game.Id), 0x00, 0x00, 0x00,
    }

    log.Println("Sending cancel for game:", game.Name)
    _, err := conn.Write(data)
    if err != nil {
        log.Println("Unable to send cancel:", err)
    }
}

func SendAnnounce(conn *net.UDPConn, game *GameInfo) {
    // FIXME int32 like a byte
    players := game.Slots - game.PlayerSlots + game.CurrentPlayers
    data := []byte {
        0xf7, 0x32, 0x10, 0x00,
        byte(game.Id),
        0x00, 0x00, 0x00,
        byte(players),
        0, 0, 0,
        byte(game.Slots),
        0, 0, 0,
    }

    log.Println("Sending announce for game:", game.Name)
    _, err := conn.Write(data)
    if err != nil {
        log.Println("Unable to send announce:", err)
    }
}

func Browse(local *net.UDPConn, remoteConn *net.UDPConn, remote *net.UDPAddr, clientVersion ClientVersion) {
    var game *GameInfo = nil
    timepoint := time.Now()
    updateGameInfo := func(g *GameInfo) {
        if g == nil {
            if game != nil {
                if time.Now().After(timepoint.Add(3 * time.Second)) {
                    SendCancel(local, game)
                    game = nil
                }
            }
        } else {
            game = g
            timepoint = time.Now()
            log.Println("Found game:", game)
            SendAnnounce(local, game)
        }
    }

    const (
        OK = 0
        TIMEOUT = 1
        ERROR = 2
    )

    data := make([]byte, 4096)
    readResponse := func() int {
        log.Println("Waiting for response...")
        remoteConn.SetReadDeadline(time.Now().Add(time.Second))
        n, src, err := remoteConn.ReadFromUDP(data)
        if err != nil {
            netErr, ok := err.(net.Error)
            timeout := ok && netErr.Timeout()
            if !timeout {
                log.Println("Unable to read game info:", err)
            }
            updateGameInfo(nil)
            if timeout {
                return TIMEOUT
            } else {
                return ERROR
            }
        }
        response := data[:n]

        if !bytes.Equal(src.IP, remote.IP) {
            log.Println("Invalid game info source:", src)
            updateGameInfo(nil)
            return ERROR
        }

        parsedGame, err := ParseGameInfo(response)
        if err != nil {
            log.Println("Unable to parse game info:", err)
            updateGameInfo(nil)
            return ERROR
        }
        updateGameInfo(&parsedGame)

        ChangeServerPort(response, 6112)
        local.Write(response)
        return OK
    }

    for {
        s := readResponse()
        switch (s) {
        case OK:
        case TIMEOUT:
            time.Sleep(time.Second)
            SendBrowse(remoteConn, remote, clientVersion)
        case ERROR:
        }
    }
}
