package main

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/jafari-mohammad-reza/foreverstore/p2p"
	"github.com/jafari-mohammad-reza/foreverstore/store"
)

type FileServerOpts struct {
	StorageRoot       string
	Transport         p2p.Transport
	PathTransformFunc store.PathTransformFunc
	remoteNodes       []string
}
type FileServer struct {
	StorageRoot       string
	Transport         p2p.Transport
	PathTransformFunc store.PathTransformFunc
	Opts              FileServerOpts
	peerLock          sync.Mutex
	peers             map[net.Addr]p2p.Peer
	remoteNodes       []string
	quitChan          chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		StorageRoot:       opts.StorageRoot,
		Transport:         opts.Transport,
		PathTransformFunc: opts.PathTransformFunc,
		remoteNodes:       opts.remoteNodes,
		quitChan:          make(chan struct{}),
		peerLock:          sync.Mutex{},
		peers:             make(map[net.Addr]p2p.Peer),
	}
}
func (fs *FileServer) Stop() {
	close(fs.quitChan)
}
func (fs *FileServer) loop() {
	defer fs.Transport.Close()
	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitChan:
			return
		}
	}
}
func (fs *FileServer) onPeer(p p2p.Peer) error {
	fs.peerLock.Lock()
	fs.peers[p.RemoteAddr()] = p
	defer fs.peerLock.Unlock()
	return nil
}

func (fs *FileServer) InitNetwork() error {
	slog.Info("InitNetwork", "nodes", fs.remoteNodes)
	for _, addr := range fs.remoteNodes {
		go func() {
			if err := fs.Transport.Dial(addr); err != nil {
				slog.Error("InitNetwork", "Error", err.Error())
				return
			}
		}()
	}
	return nil
}
func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}
	go fs.loop()
	if len(fs.remoteNodes) > 0 {
		return fs.InitNetwork()
	}
	return nil
}
