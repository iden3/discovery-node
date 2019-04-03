package node

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/discovery-node/discovery"
	"github.com/iden3/discovery-node/utils"
	"github.com/stretchr/testify/assert"
)

const privK0Hex = "9032531ad8736ff01515faecbd70f453ff1cc907cab51f6ce38985525d721ba1"

func privKHexToKeys(privKHex string) (*ecdsa.PrivateKey, discovery.PubK, common.Address) {
	privK, err := crypto.HexToECDSA(privKHex)
	if err != nil {
		panic(err)
	}
	pubKCrypto := privK.Public()
	pubK := discovery.PubK{*pubKCrypto.(*ecdsa.PublicKey)}
	addr := crypto.PubkeyToAddress(pubK.PublicKey)
	return privK, pubK, addr
}

func newTestKeyStorageAndAccount(privKHex, path, password string) (*keystore.KeyStore, *accounts.Account, error) {
	privK0, _, _ := privKHexToKeys(privKHex)

	// create new keystore with the privK, and new account
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.ImportECDSA(privK0, password)
	if err != nil {
		return nil, nil, err
	}
	return ks, &acc, nil
}

func TestSignAndVerify(t *testing.T) {
	password := "testpassword"
	ks, acc, err := newTestKeyStorageAndAccount(privK0Hex, "../tmp/test", password)
	assert.Nil(t, err)
	node := &NodeSrv{
		ks:  ks,
		acc: *acc,
	}
	node.ks.Unlock(node.acc, password)
	msg := []byte("test")
	sig, err := node.SignBytes(msg)
	assert.Nil(t, err)

	_, _, addr0 := privKHexToKeys(privK0Hex)
	verified := utils.VerifySignature(addr0, sig, msg)
	assert.True(t, verified)
}
