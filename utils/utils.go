package utils

import (
	"crypto/ecdsa"
	"log"

	gocrypto "crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// PublicKeyToString parses a public key to string
func PublicKeyToString(publicKey gocrypto.PublicKey) string {
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:]
}

// HashBytes performs a Keccak256 hash over the bytes
func HashBytes(b ...[]byte) []byte {
	h := crypto.Keccak256(b...)
	return h
}
