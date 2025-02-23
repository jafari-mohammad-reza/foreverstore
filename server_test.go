package main

import (
	"testing"

	"github.com/jafari-mohammad-reza/foreverstore/p2p"
	"github.com/jafari-mohammad-reza/foreverstore/store"
)

func makeServer(opts FileServerOpts) *FileServer {
	return NewFileServer(opts)
}

func TestServer(t *testing.T) {
	s1 := makeServer(FileServerOpts{
		Transport:         p2p.NewTCPTransport(p2p.NewTcpTransformOpts("3000", p2p.NewBuffDecoder())),
		StorageRoot:       "3000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
	})
	s2 := makeServer(FileServerOpts{
		Transport:         p2p.NewTCPTransport(p2p.NewTcpTransformOpts("4000", p2p.NewBuffDecoder())),
		StorageRoot:       "4000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
		remoteNodes:       []string{"3000", "5000"},
	})
	err := s1.Start()
	if err != nil {
		t.Error(err)
	}
	err = s2.Start()
	if err != nil {
		t.Error(err)
	}
}
