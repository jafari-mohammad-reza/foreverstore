package p2p

import "net"

type RpcData struct {
	From    net.Addr
	Payload []byte
}

type Transport interface {
	ListenAndAccept() error
	Close() error
	Consume() <-chan RpcData
}
