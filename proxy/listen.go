package proxy

import (
    "net"
)

func browse(local *net.UDPAddr,
            remoteListen *net.UDPAddr,
            remote *net.UDPAddr,
            clientVersion ClientVersion) {
    localConn, err := net.DialUDP("udp", nil, local)
    if err != nil {
        panic(err)
    }
    defer localConn.Close()

    remoteConn, err := net.ListenUDP("udp", remoteListen)
    if err != nil {
        panic(err)
    }
    defer remoteConn.Close()

    Browse(localConn, remoteConn, remote, clientVersion)
}

func Listen(remote string, clientVersion ClientVersion) {
    remote_udp_addr, err := net.ResolveUDPAddr("udp", remote)
    if err != nil {
        panic(err)
    }

    remote_listen_udp_addr := &net.UDPAddr{
        IP: net.IPv4(0, 0, 0, 0),
        Port: 6112,
    }

    local_udp_addr := &net.UDPAddr{
        IP: net.IPv4(255, 255, 255, 255),
        Port: 6112,
    }

    browse(local_udp_addr, remote_listen_udp_addr, remote_udp_addr, clientVersion)
}
