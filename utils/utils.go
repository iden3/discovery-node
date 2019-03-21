package utils

import (
	"crypto/ecdsa"
	"log"

	gocrypto "crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func PublicKeyToString(publicKey gocrypto.PublicKey) string {
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:]
}
