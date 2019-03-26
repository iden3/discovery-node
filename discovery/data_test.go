package discovery

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestIdParser(t *testing.T) {
	id := &Id{
		IdAddr:   common.Address{},
		Services: []Service{},
	}
	idBytes, err := id.Bytes()
	assert.Nil(t, err)
	id2, err := IdFromBytes(idBytes)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
}
func TestQueryAnswer(t *testing.T) {
	q := &Query{
		About:     common.Address{},
		From:      common.Address{},
		InfoFrom:  []byte{},
		Nonce:     0,
		PoW:       [32]byte{},
		Signature: []byte{},
	}
	qBytes, err := q.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, qBytes[:PREFIXLENGTH], QUERYMSG)
	qFromBytes, err := QueryFromBytes(qBytes)
	assert.Nil(t, err)
	assert.Equal(t, q, qFromBytes)

	a := &Answer{
		About:     common.Address{},
		From:      common.Address{},
		AgentId:   Service{},
		Services:  []Service{},
		Signature: []byte{},
	}
	aBytes, err := a.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, aBytes[:PREFIXLENGTH], ANSWERMSG)
	aFromBytes, err := AnswerFromBytes(aBytes)
	assert.Nil(t, err)
	assert.Equal(t, a, aFromBytes)

}
