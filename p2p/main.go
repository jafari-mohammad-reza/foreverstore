package p2p

import "net"

type RpcData struct {
	From    net.Addr
	Payload []byte
}

type Peer interface {
	net.Conn
	RemoteAddr() net.Addr
	Close() error
}
type Transport interface {
	ListenAndAccept() error
	Close() error
	Consume() <-chan RpcData
	Dial(string) error
}
type OnPeerMethod func(p Peer) error
