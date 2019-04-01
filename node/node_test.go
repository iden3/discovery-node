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

const privK0_hex = "9032531ad8736ff01515faecbd70f453ff1cc907cab51f6ce38985525d721ba1"

func privKHexToKeys(privK_hex string) (*ecdsa.PrivateKey, discovery.PubK, common.Address) {
	privK, err := crypto.HexToECDSA(privK_hex)
	if err != nil {
		panic(err)
	}
	pubK_crypto := privK.Public()
	pubK := discovery.PubK{*pubK_crypto.(*ecdsa.PublicKey)}
	addr := crypto.PubkeyToAddress(pubK.PublicKey)
	return privK, pubK, addr
}

func newTestKeyStorageAndAccount(path, password string) (*keystore.KeyStore, *accounts.Account, error) {
	privK0, _, _ := privKHexToKeys(privK0_hex)

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
	ks, acc, err := newTestKeyStorageAndAccount("../tmp/test", password)
	assert.Nil(t, err)
	node := &NodeSrv{
		ks:  ks,
		acc: *acc,
	}
	node.ks.Unlock(node.acc, password)
	msg := []byte("test")
	sig, err := node.SignBytes(msg)
	assert.Nil(t, err)

	_, _, addr0 := privKHexToKeys(privK0_hex)
	verified := utils.VerifySignature(addr0, sig, msg)
	assert.True(t, verified)
}
