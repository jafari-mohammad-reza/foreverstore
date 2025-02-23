package main

import (
	"github.com/jafari-mohammad-reza/foreverstore/p2p"
	"github.com/jafari-mohammad-reza/foreverstore/store"
)

func main() {
	transportOpts := p2p.NewTcpTransformOpts("3000", p2p.NewBuffDecoder())
	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_storage",
		Transport:         p2p.NewTCPTransport(transportOpts),
		PathTransformFunc: store.HashPathTransformFunc,
	}
	fileServer := NewFileServer(fileServerOpts)
	fileServer.Start()
	select {}
}
