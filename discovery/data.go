package discovery

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/discovery-node/utils"
)

// types of services
const RELAYTYPE = "relay"
const NAMESERVERTYPE = "nameserver"
const NOTIFICATIONSSERVERTYPE = "notificationsserver"
const DISCOVERYTYPE = "discovery"

// types of data packets
const PREFIXLENGTH = 7

var QUERYMSG = utils.HashBytes([]byte("querymsg"))[:PREFIXLENGTH]
var ANSWERMSG = utils.HashBytes([]byte("answermsg"))[:PREFIXLENGTH]

// Service holds the data about a node service (can be a Relay, a NameServer, a DiscoveryNode, etc)
type Service struct {
	IdAddr      common.Address
	PubK        *ecdsa.PublicKey // Public Key of the node, to receive encrypted data packets
	Url         string
	Type        string
	Mode        string // Active or Passive
	ProofServer []byte // TODO ProofClaimServer data type (to be defined)
}

// Id holds the data related to an identity
type Id struct {
	IdAddr   common.Address
	Services []Service
}

// Bytes parses the Id to byte array
func (id *Id) Bytes() ([]byte, error) {
	// maybe in the future write a byte parser&unparser
	// var b bytes.Buffer
	// b.Write((id.IdAddr.Bytes()[:]))
	// for _, s := range id.Services {
	//         b.Write()
	// }

	return json.Marshal(id)
}

// IdFromBytes parses Id data structure from a byte array
func IdFromBytes(b []byte) (*Id, error) {
	var id *Id
	err := json.Unmarshal(b, &id)
	return id, err
}

// Query is the data packet that a node sends to discover data about one identity
type Query struct {
	About     common.Address // About Who
	From      common.Address
	InfoFrom  []byte // TODO to be defined
	Timestamp int64
	Nonce     uint64
	PoW       [32]byte // TODO for the moment Keccak256
	Signature []byte
}

// Bytes parses the Query to byte array
func (q *Query) Bytes() ([]byte, error) {
	b, err := json.Marshal(q)
	if err != nil {
		return b, err
	}
	var r []byte
	r = append(r, QUERYMSG...)
	r = append(r, b...)
	return r, nil
}

// QureyFromBytes parses Query data structure from a byte array
func QueryFromBytes(b []byte) (*Query, error) {
	if !bytes.Equal(b[:PREFIXLENGTH], QUERYMSG) {
		return nil, errors.New("Not query type")
	}
	var q *Query
	err := json.Unmarshal(b[PREFIXLENGTH:], &q)
	return q, err
}

// Answer is the data packet that a node sends when answering to a Query data packet
type Answer struct {
	About     common.Address
	From      common.Address
	AgentId   Service
	Services  []Service
	Timestamp int64
	Signature []byte
}

// Bytes parses the Answer to byte array
func (a *Answer) Id() *Id {
	return &Id{
		IdAddr:   a.About,
		Services: a.Services,
	}
}

// Bytes parses the Answer to byte array
func (a *Answer) Bytes() ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return b, err
	}
	var r []byte
	r = append(r, ANSWERMSG...)
	r = append(r, b...)
	return r, nil
}

// AnswerFromBytes parses Answer data structure from a byte array
func AnswerFromBytes(b []byte) (*Answer, error) {
	if !bytes.Equal(b[:PREFIXLENGTH], ANSWERMSG) {
		return nil, errors.New("Not answer type")
	}
	var a *Answer
	err := json.Unmarshal(b[PREFIXLENGTH:], &a)
	return a, err
}
