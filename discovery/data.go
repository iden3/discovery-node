package discovery


// types of services
const RELAYTYPE = "relay"
const NAMESERVERTYPE = "nameserver"
const NOTIFICATIONSSERVERTYPE = "notificationsserver"
const DISCOVERYTYPE = "discovery"

// Query is the data packet that a node sends to discover data about one identity
type Query struct {
	About string // About Who
	From common.Address
	InfoFrom []byte // TODO to be defined
	Nonce uint64
	PoW [32]byte // TODO for the moment Keccak256
	Signature []byte
}

// Answer is the data packet that a node sends when answering to a Query data packet
type Answer struct {
	About string
	From string
	AgentId Service
	Services []Service
	Signature []byte
}

// Service holds the data about a node service (can be a Relay, a NameServer, a DiscoveryNode, etc)
type Service struct {
	IdAddr common.Address
	PubK []byte // Public Key of the node, to receive encrypted data packets
	Url string
	Type string
	Mode string // Active or Passive
	ProofServer []byte // TODO ProofClaimServer data type (to be defined)
}
