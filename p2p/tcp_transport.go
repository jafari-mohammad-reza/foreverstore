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
	peers  map[net.Addr]*TCPPeer
	mu     sync.RWMutex
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	lnAddress := ":" + listenAddress

	return &TCPTransport{
		lnAddr: lnAddress,
		peers:  make(map[net.Addr]*TCPPeer),
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
	peer := NewTCPPeer(conn, true)
	t.peers[peer.conn.LocalAddr()] = peer
}
