package p2p

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTrans(t *testing.T) {
	listenAddr := "4000"
	tr := NewTCPTransport(listenAddr)
	assert.Equal(t, tr.lnAddr, ":"+listenAddr)
	assert.Nil(t, tr.ListenAndAccept())
	assert.Len(t, tr.peers, 0)
	conn, err := net.Dial("tcp", tr.lnAddr)
	conn.Write([]byte("test peer"))
	assert.Nil(t, err)
	assert.Len(t, tr.peers, 1)
}
