package main

import (
	"github.com/jafari-mohammad-reza/foreverstore/p2p"
)

func main() {
	tr := p2p.NewTCPTransport("4000")
	tr.ListenAndAccept()
	select {}
}
