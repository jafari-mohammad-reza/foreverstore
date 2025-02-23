package main

import (
	"testing"
	"time"

	"github.com/jafari-mohammad-reza/foreverstore/p2p"
	"github.com/jafari-mohammad-reza/foreverstore/store"
	"github.com/stretchr/testify/assert"
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
		remoteNodes:       []string{"3000"},
	})
	s3 := makeServer(FileServerOpts{
		Transport:         p2p.NewTCPTransport(p2p.NewTcpTransformOpts("5000", p2p.NewBuffDecoder())),
		StorageRoot:       "5000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
		remoteNodes:       []string{"3000", "4000", "6000"},
	})
	err := s1.Start()
	if err != nil {
		t.Error(err)
	}
	err = s2.Start()
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	err = s3.Start() // should return dial error
	if err != nil {
		assert.NotNil(t, err)
	}
}
