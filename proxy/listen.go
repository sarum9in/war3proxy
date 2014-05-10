package proxy

import (
    "net"
)

func handle(local *net.TCPConn, host *net.TCPAddr) {
    defer local.Close()

    remote, err := net.DialTCP("tcp", nil, host)
    if err != nil {
        panic(err)
    }
    defer remote.Close()

    Proxy(local, remote)
}

func Listen(host string, port uint) {
    host_addr, err := net.ResolveTCPAddr("tcp", host + string(port))
    if err != nil {
        panic(err)
    }

    local_addr, err := net.ResolveTCPAddr("tcp", "localhost:6112")
    if err != nil {
        panic(err)
    }

    listener, err := net.ListenTCP("tcp", local_addr)
    if err != nil {
        panic(err)
    }

    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            panic(err)
        }
        go handle(conn, host_addr)
    }
}
