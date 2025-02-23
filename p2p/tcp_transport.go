package p2p

import (
	"errors"
	"io"
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
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

type TransportOpts struct {
	lnAddr   string
	ln       net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
	peerLock sync.RWMutex
	rpcChan  chan RpcData
	decoder  Decoder
}

func NewTcpTransformOpts(lnAddress string, decoder Decoder) TransportOpts {
	return TransportOpts{
		lnAddr:   lnAddress,
		mu:       sync.RWMutex{},
		rpcChan:  make(chan RpcData),
		peers:    make(map[net.Addr]Peer),
		peerLock: sync.RWMutex{},
		decoder:  decoder,
	}
}

type TCPTransport struct {
	TransportOpts
}

func NewTCPTransport(opts TransportOpts) *TCPTransport {
	opts.lnAddr = ":" + opts.lnAddr

	return &TCPTransport{
		TransportOpts: opts,
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
	slog.Info("TcpTransport ListenAndAccept", "Port", t.lnAddr)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			slog.Error("acceptLoop", "Error", err.Error())
		}
		go t.handleConn(conn, false)
	}
}
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	t.peerLock.Lock()
	t.peers[conn.RemoteAddr()] = NewTCPPeer(conn, outbound)
	t.peerLock.Unlock()
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error("TcpTransport handleConn", "Error", err.Error())
		}
	}()
	buff := make([]byte, 1028)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Client disconnected", "addr", conn.RemoteAddr())
			} else {
				slog.Error("handleConn read error", "addr", conn.RemoteAddr(), "error", err)
			}
			break
		}
		select {
		case t.rpcChan <- RpcData{From: conn.RemoteAddr(), Payload: buff[:n]}:
		default:
			slog.Warn("Dropping message: channel full", "addr", conn.RemoteAddr())
		}
	}
}
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", ":"+addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}
func (t *TCPTransport) Consume() <-chan RpcData {
	return t.rpcChan
}
func (t *TCPTransport) Close() error {
	return t.ln.Close()
}
