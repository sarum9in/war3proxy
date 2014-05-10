package main

import (
    "flag"
    "fmt"
    "./proxy"
)

var host string
var port uint

func init() {
    flag.StringVar(&host, "host", "localhost", "address (or IP) of host")
    flag.UintVar(&port, "port", 6112, "port of host (if not standard 6112)")
}

func main() {
    flag.Parse()
    fmt.Printf("Attempt to connect to %s:%d\n", host, port)
    proxy.Listen(host, port)
}
