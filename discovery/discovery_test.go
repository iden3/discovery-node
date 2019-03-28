package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiscoveryService(t *testing.T) {
	_, err := NewDiscoveryService(addr0, []byte{}, &pubK0.PublicKey, "url", "Active", []byte{})
	assert.Nil(t, err)
}
func TestNewQueryPacket(t *testing.T) {
	dscsrv, err := NewDiscoveryService(addr0, []byte{}, &pubK0.PublicKey, "url", "Active", []byte{})
	assert.Nil(t, err)
	q, err := dscsrv.NewQueryPacket(addr1)
	assert.Nil(t, err)

	qBytes, err := q.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, qBytes[:PREFIXLENGTH], QUERYMSG)
	qFromBytes, err := QueryFromBytes(qBytes)
	assert.Nil(t, err)
	assert.Equal(t, q, qFromBytes)
}
func TestNewAnswerPacket(t *testing.T) {
	dscsrv, err := NewDiscoveryService(addr0, []byte{}, &pubK0.PublicKey, "url", "Active", []byte{})
	assert.Nil(t, err)
	query, err := dscsrv.NewQueryPacket(addr1)
	assert.Nil(t, err)

	id := &Id{
		IdAddr:   addr1,
		Services: []Service{},
	}
	_, err = dscsrv.NewAnswerPacket(query, id)
	assert.Nil(t, err)

	// check that if the id.IdAddr is different than query.AboutId, gives an error
	id.IdAddr = addr2
	_, err = dscsrv.NewAnswerPacket(query, id)
	assert.NotNil(t, err)

}
