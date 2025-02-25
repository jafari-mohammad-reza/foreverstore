package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/jafari-mohammad-reza/foreverstore/p2p"
	"github.com/jafari-mohammad-reza/foreverstore/store"
	"github.com/stretchr/testify/assert"
)

func makeServer(opts FileServerOpts, tcpOpts p2p.TransportOpts) *FileServer {
	tcpTransport := p2p.NewTCPTransport(tcpOpts)
	server := NewFileServer(opts)
	tcpTransport.OnPeer = server.onPeer
	server.Transport = tcpTransport
	return server
}
func TestServer(t *testing.T) {
	s1 := makeServer(FileServerOpts{
		StorageRoot:       "3000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
	}, p2p.NewTcpTransformOpts("3000", p2p.NewBuffDecoder()))
	s2 := makeServer(FileServerOpts{
		StorageRoot:       "4000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
		remoteNodes:       []string{"3000"},
	}, p2p.NewTcpTransformOpts("4000", p2p.NewBuffDecoder()))
	s3 := makeServer(FileServerOpts{
		StorageRoot:       "5000_storage",
		PathTransformFunc: store.HashPathTransformFunc,
		remoteNodes:       []string{"3000", "4000", "6000"},
	}, p2p.NewTcpTransformOpts("5000", p2p.NewBuffDecoder()))
	err := s1.Start()
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	err = s2.Start()
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	err = s3.Start() // should return dial error
	if err != nil {
		assert.NotNil(t, err)
	}
	time.Sleep(time.Second)

	writtenData := bytes.NewReader([]byte("test data"))
	err = s2.storeData("storeDataTestKey", writtenData)
	if err != nil {
		assert.Nil(t, err)
	}
	time.Sleep(time.Second)
	defer func() {
		os.RemoveAll("3000_storage")
		os.RemoveAll("4000_storage")
		os.RemoveAll("5000_storage")
	}()
}
