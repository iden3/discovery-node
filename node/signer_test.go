package node

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const privK1Hex = "36aa8fc62a7e482a5f0bfb0c054601fdf2d1cb1d83ff9a1915c53f25ac7968c3"

func TestSignBytes(t *testing.T) {
	password := "testpassword"
	ks, acc, err := newTestKeyStorageAndAccount(privK1Hex, "../tmp/test"+time.Now().String(), password)
	assert.Nil(t, err)
	node := &NodeSrv{
		ks:  ks,
		acc: *acc,
	}
	node.ks.Unlock(node.acc, password)
	msg := []byte("test")
	sig, err := node.SignBytes(msg)
	assert.Nil(t, err)

	assert.Equal(t, hex.EncodeToString(sig), "e07438d24681aae2e49c0e57562d0671b54a459e238027d09a59e508ea65b32e603eafb4271ce4773c56c42cebfe7b3f62d598dac0debb7098b66e435d34aa951c")
}
