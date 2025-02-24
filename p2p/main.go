package p2p

import "net"

type RpcData struct {
	From    net.Addr
	Payload []byte
}

type Peer interface {
	RemoteAddr() net.Addr
}
type Transport interface {
	ListenAndAccept() error
	Close() error
	Consume() <-chan RpcData
	Dial(string) error
}
type OnPeerMethod func(p Peer) error
