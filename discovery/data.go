package discovery

import (
	"encoding/json"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)


// types of services
const RELAYTYPE = "relay"
const NAMESERVERTYPE = "nameserver"
const NOTIFICATIONSSERVERTYPE = "notificationsserver"
const DISCOVERYTYPE = "discovery"

// Service holds the data about a node service (can be a Relay, a NameServer, a DiscoveryNode, etc)
type Service struct {
	IdAddr common.Address
	PubK *ecdsa.PublicKey // Public Key of the node, to receive encrypted data packets
	Url string
	Type string
	Mode string // Active or Passive
	ProofServer []byte // TODO ProofClaimServer data type (to be defined)
}

// Id holds the data related to an identity
type Id struct {
	IdAddr common.Address
	Services []Service
}

func (id *Id) Bytes() ([]byte, error) {
	// maybe in the future write a byte parser&unparser
	// var b bytes.Buffer
	// b.Write((id.IdAddr.Bytes()[:]))
	// for _, s := range id.Services {
	//         b.Write()
	// }

	return json.Marshal(id)
}
func IdFromBytes(b []byte) (*Id, error) {
	var id *Id
	err := json.Unmarshal(b, &id)
	return id, err
}

// Query is the data packet that a node sends to discover data about one identity
type Query struct {
	About common.Address // About Who
	From common.Address
	InfoFrom []byte // TODO to be defined
	Nonce uint64
	PoW [32]byte // TODO for the moment Keccak256
	Signature []byte
}

// Answer is the data packet that a node sends when answering to a Query data packet
type Answer struct {
	About common.Address
	From common.Address
	AgentId Service
	Services []Service
	Signature []byte
}

