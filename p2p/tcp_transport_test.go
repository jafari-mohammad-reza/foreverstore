package p2p

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpTrans(t *testing.T) {
	opts := NewTcpTransformOpts("4000", NewBuffDecoder())

	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.lnAddr, ":"+opts.lnAddr)
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
