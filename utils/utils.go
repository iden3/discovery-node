package utils

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"

	gocrypto "crypto"

	"github.com/ethereum/go-ethereum/common"
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

// EthHash performs a Keccak256 hash of the EthPrefix + byte array
func EthHash(b []byte) []byte {
	header := fmt.Sprintf("%s%d", "\x19Ethereum Signed Message:\n", len(b))
	return HashBytes([]byte(header), b)
}

// VerifySignature verifies that the signature of the msg is made by the private key of the given address
func VerifySignature(addr common.Address, sig, msg []byte) bool {
	h := EthHash(msg)
	sig[64] -= 27

	recoveredPub, err := crypto.Ecrecover(h, sig)
	if err != nil {
		return false
	}
	pubK, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*pubK)
	return bytes.Equal(addr.Bytes(), recoveredAddr.Bytes())

}
