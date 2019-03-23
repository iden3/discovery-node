package discovery

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIdParser(t *testing.T) {
	id := &Id{
		IdAddr: common.Address{},
		Services: []Service{},
	}
	idBytes, err := id.Bytes()
	assert.Nil(t, err)
	id2, err := IdFromBytes(idBytes)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
}
