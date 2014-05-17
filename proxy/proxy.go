package proxy

import (
    "log"
    "net"
    "time"
    "../warcraft"
)

// Send browse packet to remote via conn for clientVersion.
// This is used to discover remote's games.
func SendBrowsePacket(conn *net.UDPConn, remote *net.UDPAddr, clientVersion warcraft.ClientVersion) {
    browse := warcraft.NewBrowsePacket(clientVersion)
    log.Printf("Sending browse packet for client version: %q\n", clientVersion)

    _, err := conn.WriteTo(browse.Bytes(), remote)
    if err != nil {
        log.Println("Unable to send browse:", err)
    }
}

// Send cancel packet via conn for game.
// This is used to cancel previously announced game.
func SendCancelPacket(conn *net.UDPConn, game *warcraft.GameInfo) {
    cancel := warcraft.NewCancelPacket(game.Id)

    log.Printf("Sending cancel packet for game: %q\n", game.Name)
    _, err := conn.Write(cancel.Bytes())
    if err != nil {
        log.Println("Unable to send cancel:", err)
    }
}

// Send announce packet via conn for game.
// This is used to make game known for clients
// or update already known game's info.
func SendAnnouncePacket(conn *net.UDPConn, game *warcraft.GameInfo) {
    players := game.Slots - game.PlayerSlots + game.CurrentPlayers
    announce := warcraft.NewAnnouncePacket(game.Id, players, game.Slots)

    log.Printf("Sending announce packet for game: %q\n", game.Name)
    _, err := conn.Write(announce.Bytes())
    if err != nil {
        log.Println("Unable to send announce:", err)
    }
}

// If there is no response from remote server
// for GameTimeout time the game is considered lost.
var GameTimeout = 3 * time.Second

// Pass-through remote client to local network.
func Browse(local *net.UDPConn,
            remoteConn *net.UDPConn,
            remote *net.UDPAddr,
            clientVersion warcraft.ClientVersion) {
    var game *warcraft.GameInfo = nil
    timepoint := time.Now()

    // update game sending announce or cancel packet
    updateGameInfo := func(g *warcraft.GameInfo) {
        if g == nil {
            if game != nil {
                // remote client does not send cancel packets,
                // so it is necessary to track game's existence
                if time.Now().After(timepoint.Add(GameTimeout)) {
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
    // read remote client's response
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

        if !src.IP.Equal(remote.IP) || src.Port != remote.Port {
            log.Println("Invalid game info source:", src.IP, "!=", remote.IP)
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

        // update remote's server port to local
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
