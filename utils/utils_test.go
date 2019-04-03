package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyToString(t *testing.T) {
	privK0Hex := "9032531ad8736ff01515faecbd70f453ff1cc907cab51f6ce38985525d721ba1"
	privK0, err := crypto.HexToECDSA(privK0Hex)
	if err != nil {
		panic(err)
	}
	pubK0_crypto := privK0.Public()
	pubK0 := pubK0_crypto.(*ecdsa.PublicKey)

	pubKStr := PublicKeyToString(pubK0)
	assert.Equal(t, pubKStr, "324b1aff3e5cc35f1b5df87f6bb89051e90a93eaffb5dc41ac12621f64cb479c9b55ce5d0776fd196dc3a9ad2461918ff4bd0fc1dec659758c75592401974915")
}

func TestHashBytes(t *testing.T) {
	msg := []byte("test")
	h := HashBytes(msg)
	assert.Equal(t, hex.EncodeToString(h), "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658")
}

func TestEthHash(t *testing.T) {
	msg := []byte("test")
	h := EthHash(msg)
	assert.Equal(t, hex.EncodeToString(h), "4a5c5d454721bbbb25540c3317521e71c373ae36458f960d2ad46ef088110e95")
}

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
