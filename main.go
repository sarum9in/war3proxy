package main

import (
    "flag"
    "fmt"
    "log"
    "./proxy"
)

var host string
var port uint
var client proxy.ClientVersion
var tft bool
var version uint

func init() {
    flag.StringVar(&host, "host", "localhost", "address (or IP) of host")
    flag.UintVar(&port, "port", 6112, "port of host (if not standard 6112)")
    flag.BoolVar(&tft, "tft", true, "Use TFT expansion")
    flag.UintVar(&version, "version", 26, "version of client (for 1.x enter x)")
}

func main() {
    flag.Parse()
    if tft {
    	client.Expansion = proxy.TftExpansion
    } else {
    	client.Expansion = proxy.RawExpansion
    }
    client.Version = byte(version)
    addr := fmt.Sprintf("%s:%d", host, port)
    log.Println("Proxying to", addr)
    proxy.Listen(addr, client)
}
