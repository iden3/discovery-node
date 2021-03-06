package discovery

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/discovery-node/utils"
)

// DISCOVERYVERSION specifies the version of the discovery protocol
const DISCOVERYVERSION = "v0.0.1"

// types of services
const RELAYTYPE = "relay"
const NAMESERVERTYPE = "nameserver"
const NOTIFICATIONSSERVERTYPE = "notificationsserver"
const DISCOVERYTYPE = "discovery"

// types of data packets
const PREFIXLENGTH = 7

// QUERYMSG is the identifier of Query messages
var QUERYMSG = utils.HashBytes([]byte("querymsg"))[:PREFIXLENGTH]

// ANSWERMSG is the identifier of Answer messages
var ANSWERMSG = utils.HashBytes([]byte("answermsg"))[:PREFIXLENGTH]

// PubK is a ecdsa.PublicKey with json marshal and unmarshal functions that convert the key into a byte array
type PubK struct {
	ecdsa.PublicKey
}

func (pubK *PubK) MarshalJSON() ([]byte, error) {
	return json.Marshal(crypto.CompressPubkey(&pubK.PublicKey))
}
func (pubK *PubK) UnmarshalJSON(b []byte) error {
	var bb []byte
	err := json.Unmarshal(b, &bb)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p, err := crypto.DecompressPubkey(bb)
	if err != nil {
		return err
	}
	pubK.PublicKey = *p
	return nil
}

// String returns the PubK in string format
func (pubK *PubK) String() string {
	publicKeyBytes := crypto.FromECDSAPub(&pubK.PublicKey)
	return hexutil.Encode(publicKeyBytes)[4:]
}

// Service holds the data about a node service (can be a Relay, a NameServer, a DiscoveryNode, etc)
type Service struct {
	IdAddr       common.Address
	KademliaAddr []byte // Kademlia address
	PssPubK      PubK   // Public Key of the pss node, to receive encrypted data packets
	Url          string
	Type         string // TODO define type specification (relay, nameserver, etc)
	Mode         string // Active or Passive(gateway) (this only affects to discovery-node's type)
	ProofService []byte // TODO ProofClaimService data type (to be defined)
}

// Id holds the data related to an identity
type Id struct {
	IdAddr   common.Address
	Services []Service
}

// Bytes parses the Id to byte array
func (id *Id) Bytes() ([]byte, error) {
	// maybe in the future write a byte parser&unparser
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
	Version          string         // version of the protocol
	MsgId            string         // random msg id, to identify and relate Query and Answer
	AboutId          common.Address // About Who is requesting data (about which identity address)
	RequesterId      common.Address
	RequesterKAddr   []byte // Kademlia address
	RequesterPssPubK PubK   // Public Key of the pss node requester, to receive encrypted data packets
	InfoFrom         []byte // TODO to be defined
	Timestamp        int64
	Nonce            uint64 // for the PoW
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

// QueryFromBytes parses Query data structure from a byte array
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
	Version   string // version of the protocol
	MsgId     string // random msg id, to identify and relate Query and Answer
	AboutId   common.Address
	FromId    common.Address
	AgentId   Service
	Services  []Service
	Timestamp int64
	Signature []byte
}

// Id returns a pointer to an Id object from the Answer data packet
func (a *Answer) Id() *Id {
	return &Id{
		IdAddr:   a.AboutId,
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

// Copy returns a pointer to a copy of the Answer data packet
func (a *Answer) Copy() *Answer {
	return &Answer{
		Version:   a.Version,
		MsgId:     a.MsgId,
		AboutId:   a.AboutId,
		FromId:    a.FromId,
		AgentId:   a.AgentId,
		Services:  a.Services,
		Timestamp: a.Timestamp,
		Signature: a.Signature,
	}
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
