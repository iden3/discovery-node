package utils

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {
	sigHex := "d192babb565316dfb5905c00043bbba84d6066ce0d24e0109ce5b63dd453eec1484b89c0b417fa13bdf6607e615875c4bb0c5a6014bfd47c3c1216ebc41d96c81c"
	addrHex := "0x47536b157a638db4d88b6c71f1143d66611a5d19"
	msg := []byte("test")

	sig, err := hex.DecodeString(sigHex)
	assert.Nil(t, err)
	addr := common.HexToAddress(addrHex)

	verified := VerifySignature(addr, sig, msg)
	assert.True(t, verified)
}
