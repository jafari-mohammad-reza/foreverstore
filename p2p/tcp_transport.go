package p2p

import (
	"log/slog"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool // case we dial the connection not accept it
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	lnAddr string
	ln     net.Listener

	mu      sync.RWMutex
	rpcChan chan RpcData
	decoder Decoder
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	lnAddress := ":" + listenAddress

	return &TCPTransport{
		lnAddr:  lnAddress,
		mu:      sync.RWMutex{},
		rpcChan: make(chan RpcData),
		decoder: NewBuffDecoder(),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.lnAddr)
	if err != nil {
		slog.Error("listenAndAccept", "Error", err.Error())
		return err
	}
	t.ln = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			slog.Error("acceptLoop", "Error", err.Error())
		}
		go t.handleConn(conn)
	}
}
func (t *TCPTransport) handleConn(conn net.Conn) {

	rpcData := RpcData{
		conn.RemoteAddr(),
		[]byte{},
	}
	for {
		buff := make([]byte, 1028)
		n, err := conn.Read(buff)
		if err != nil {
			slog.Error("handleConn read", "Error", err.Error())
			conn.Close()
			break
		}
		rpcData.Payload = buff[:n]
		t.rpcChan <- rpcData
	}
}
func (t *TCPTransport) Consume() <-chan RpcData {
	return t.rpcChan
}
func (t *TCPTransport) Close() error {
	return t.ln.Close()
}
