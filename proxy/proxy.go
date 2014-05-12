package proxy

import (
    "bytes"
    "log"
    "net"
    "time"
    "../warcraft"
)

func SendBrowsePacket(conn *net.UDPConn, remote *net.UDPAddr, clientVersion warcraft.ClientVersion) {
    browse := warcraft.NewBrowsePacket(clientVersion)
    log.Printf("Sending browse packet for client version: %q\n", clientVersion)

    _, err := conn.WriteTo(browse.Bytes(), remote)
    if err != nil {
        log.Println("Unable to send browse:", err)
    }
}

func SendCancelPacket(conn *net.UDPConn, game *warcraft.GameInfo) {
    cancel := warcraft.NewCancelPacket(game.Id)

    log.Printf("Sending cancel packet for game: %q\n", game.Name)
    _, err := conn.Write(cancel.Bytes())
    if err != nil {
        log.Println("Unable to send cancel:", err)
    }
}

func SendAnnouncePacket(conn *net.UDPConn, game *warcraft.GameInfo) {
    players := game.Slots - game.PlayerSlots + game.CurrentPlayers
    announce := warcraft.NewAnnouncePacket(game.Id, players, game.Slots)

    log.Printf("Sending announce packet for game: %q\n", game.Name)
    _, err := conn.Write(announce.Bytes())
    if err != nil {
        log.Println("Unable to send announce:", err)
    }
}

func Browse(local *net.UDPConn, remoteConn *net.UDPConn, remote *net.UDPAddr, clientVersion warcraft.ClientVersion) {
    var game *warcraft.GameInfo = nil
    timepoint := time.Now()
    updateGameInfo := func(g *warcraft.GameInfo) {
        if g == nil {
            if game != nil {
                if time.Now().After(timepoint.Add(3 * time.Second)) {
                    SendCancelPacket(local, game)
                    game = nil
                }
            }
        } else {
            game = g
            timepoint = time.Now()
            log.Println("Found game:", game)
            SendAnnouncePacket(local, game)
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

        parsedGame, err := warcraft.ParseGameInfo(response)
        if err != nil {
            log.Println("Unable to parse game info:", err)
            updateGameInfo(nil)
            return ERROR
        }
        updateGameInfo(&parsedGame)

        warcraft.ChangeServerPort(response, 6112)
        local.Write(response)
        return OK
    }

    for {
        s := readResponse()
        switch (s) {
        case OK:
        case TIMEOUT:
            time.Sleep(time.Second)
            SendBrowsePacket(remoteConn, remote, clientVersion)
        case ERROR:
        }
    }
}
