package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
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
	Transport   p2p.Transport
	Opts        FileServerOpts
	peerLock    sync.Mutex
	peers       map[net.Addr]p2p.Peer
	remoteNodes []string
	quitChan    chan struct{}
	store       *store.Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	store := store.NewStore(store.StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	})
	return &FileServer{
		Transport:   opts.Transport,
		remoteNodes: opts.remoteNodes,
		quitChan:    make(chan struct{}),
		peerLock:    sync.Mutex{},
		peers:       make(map[net.Addr]p2p.Peer),
		store:       store,
		Opts:        opts,
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
			fmt.Printf("Message received at %s: %+v\n", fs.Opts.StorageRoot, msg)
			var p StoreDataPayload
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("consumed payload %+v\n from %s - %+v", p, msg.From, fs.remoteNodes)
		case <-fs.quitChan:
			return
		}
	}
}
func (fs *FileServer) onPeer(p p2p.Peer) error {
	fmt.Println("server onPeer", p.RemoteAddr(), fs.Opts.StorageRoot)
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()
	fs.peers[p.RemoteAddr()] = p
	fmt.Printf("Current peers: %+v\n", fs.peers)
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

func (fs *FileServer) storeData(key string, r io.Reader) error {
	buff := new(bytes.Buffer)
	tee := io.TeeReader(r, buff)
	if err := fs.store.WriteStream(key, tee); err != nil {
		return err
	}
	payload := &StoreDataPayload{
		Key:  key,
		Data: buff.Bytes(),
	}
	return fs.broadcast(payload)
}

type StoreDataPayload struct {
	Key  string
	Data []byte
}

func (fs *FileServer) broadcast(payload *StoreDataPayload) error {
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	fmt.Printf("Broadcasting from %s payload %+v to peers: %+v\n", fs.Opts.StorageRoot, payload, fs.peers)
	peers := []io.Writer{}
	for _, p := range fs.peers {
		peers = append(peers, p)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(payload)
}
