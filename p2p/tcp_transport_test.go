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
	conn, err := net.Dial("tcp", tr.lnAddr)
	testmsg := "test message"
	conn.Write([]byte(testmsg))
	assert.Nil(t, err)
	for chanData := range tr.Consume() {
		assert.Equal(t, string(chanData.Payload), testmsg)
		break
	}
}
