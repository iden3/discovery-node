package discovery

import (
	"crypto/ecdsa"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

const privK0Hex = "9032531ad8736ff01515faecbd70f453ff1cc907cab51f6ce38985525d721ba1"
const privK1Hex = "36aa8fc62a7e482a5f0bfb0c054601fdf2d1cb1d83ff9a1915c53f25ac7968c3"
const privK2Hex = "a71f5762a9d698b88f97dfc0962fc65151dfddfdde366a1362cff0551daca797"
const privK3Hex = "e7f83fd5cdd1fa6ee0247262d5cbb61cfbba3a7300948a1a6063009b6c2b997e"

// 0 discovery-node requester keys
var privK0 *ecdsa.PrivateKey
var pubK0 PubK
var addr0 common.Address

// 1 discovery-node id_agent keys
var privK1 *ecdsa.PrivateKey
var pubK1 PubK
var addr1 common.Address

// 2 user id keys
var privK2 *ecdsa.PrivateKey
var pubK2 PubK
var addr2 common.Address

// 3 relay service keys
var privK3 *ecdsa.PrivateKey
var pubK3 PubK
var addr3 common.Address

var serviceNode0 Service
var serviceNode1 Service
var serviceRelay Service

func init() {
	var err error

	// 0
	privK0, err = crypto.HexToECDSA(privK0Hex)
	if err != nil {
		panic(err)
	}
	pubK0_crypto := privK0.Public()
	pubK0 = PubK{*pubK0_crypto.(*ecdsa.PublicKey)}
	addr0 = crypto.PubkeyToAddress(pubK0.PublicKey)

	// 1
	privK1, err = crypto.HexToECDSA(privK1Hex)
	if err != nil {
		panic(err)
	}
	pubK1_crypto := privK1.Public()
	pubK1 = PubK{*pubK1_crypto.(*ecdsa.PublicKey)}
	addr1 = crypto.PubkeyToAddress(pubK1.PublicKey)

	// 2
	privK2, err = crypto.HexToECDSA(privK2Hex)
	if err != nil {
		panic(err)
	}
	pubK2_crypto := privK2.Public()
	pubK2 = PubK{*pubK2_crypto.(*ecdsa.PublicKey)}
	addr2 = crypto.PubkeyToAddress(pubK2.PublicKey)

	// 3
	privK3, err = crypto.HexToECDSA(privK3Hex)
	if err != nil {
		panic(err)
	}
	pubK3_crypto := privK3.Public()
	pubK3 = PubK{*pubK3_crypto.(*ecdsa.PublicKey)}
	addr3 = crypto.PubkeyToAddress(pubK3.PublicKey)

	serviceNode0 = Service{
		IdAddr:       addr0,
		PssPubK:      pubK0,
		Url:          "",
		Type:         "discovery-node",
		Mode:         "ACTIVE",
		ProofService: []byte{},
	}
	serviceNode1 = Service{
		IdAddr:       addr1,
		PssPubK:      pubK1,
		Url:          "",
		Type:         "discovery-node",
		Mode:         "ACTIVE",
		ProofService: []byte{},
	}
	serviceRelay = Service{
		IdAddr:       addr3,
		PssPubK:      pubK3,
		Url:          "https://relay.domain.eth",
		Type:         "iden3-relay",
		Mode:         "",
		ProofService: []byte{},
	}
}

func TestPubK(t *testing.T) {
	mPubK0, err := json.Marshal(&pubK0)
	assert.Nil(t, err)
	var p PubK
	err = json.Unmarshal(mPubK0, &p)
	assert.Nil(t, err)
	assert.Equal(t, pubK0, p)
}

func TestIdParser(t *testing.T) {
	id := &Id{
		IdAddr:   addr2,
		Services: []Service{serviceRelay},
	}
	idBytes, err := id.Bytes()
	assert.Nil(t, err)
	id2, err := IdFromBytes(idBytes)
	assert.Nil(t, err)
	assert.Equal(t, id, id2)
}

func TestQueryAnswer(t *testing.T) {

	q := &Query{
		Version:          "v0.0.1",
		AboutId:          addr2,
		RequesterId:      addr0,
		RequesterPssPubK: pubK0,
		InfoFrom:         []byte{},
		Nonce:            0,
	}
	qBytes, err := q.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, qBytes[:PREFIXLENGTH], QUERYMSG)
	qFromBytes, err := QueryFromBytes(qBytes)
	assert.Nil(t, err)
	assert.Equal(t, q, qFromBytes)

	a := &Answer{
		AboutId:   addr2,
		FromId:    addr1,
		AgentId:   serviceNode1,
		Services:  []Service{serviceRelay},
		Signature: []byte{},
	}
	aBytes, err := a.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, aBytes[:PREFIXLENGTH], ANSWERMSG)
	aFromBytes, err := AnswerFromBytes(aBytes)
	assert.Nil(t, err)
	assert.Equal(t, a, aFromBytes)

}
