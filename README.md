# discovery-node
Draft implementation of `discovery-node` of the decentralized discovery protocol over Pss Swarm


## Overview

![network00](https://raw.githubusercontent.com/iden3/discovery-node/master/docs/network00.png "network00")

Types of node:
	- `passive`: are the nodes that only perform petitions, acting as gateways to the discovery network
	- `active`: are the nodes that are answering requests

Entities:
- `Issuer`: entity that wants to know about one identity
- `ID_AGENT`/`Discovery Server`: entity that knows the info about the identity, and is listening in `Swarm Pss` in the topic `id_discovery`

#### Sample flow
- `discovery-node` receives an http petition asking for an identity info, from now, the `discovery-node` will be the `Issuer`
- `Issuer` creates `Query` packet asking for who is the relay of identity `john@domain.eth`
- `Issuer` sends the `Query` packet into the `Swarm Pss` network under the topic `id_discovery`
- the `ID_AGENT` server of that identity will receive the `Query` packet and will see that is a user under its umbrella
- `ID_AGENT` server will answer the `Answer` packet (with the proofs of validity) to the `Issuer`
- `Issuer` receives the `Answer` packet, and now knows how to reach the Relay node of `john@domain.eth`

#### Node Storage
The `discovery-node` data storage is a leveldb database. It's organized with prefixes, where each type of data is stored under a prefix.

Databases:
- `dbOwnIds`: holds the data about the identities that the `discovery-node` manages
- `dbAnswCache`: holds the data about the discovered identites. Each data packet of a discovered identity, has a `timestamp`, the data packets are valid under a time window where the `timestamp` allows to determine if it's already valed or is too old

#### Data structures

```go
// Service holds the data about a node service (can be a Relay, a NameServer, a DiscoveryNode, etc)
type Service struct {
	IdAddr      common.Address
	PubK        *ecdsa.PublicKey // (optional ) Public Key of the node, to receive encrypted data packets
	Url         string
	Type        string
	Mode        string // Active or Passive
	ProofServer []byte // ProofClaimServer
}

// Query is the data packet that a node sends to discover data about one identity
type Query struct {
	About     common.Address // About Who
	From      common.Address
	InfoFrom  []byte // to be defined
	Timestamp int64
	Nonce     uint64
	PoW       [32]byte // for the moment Keccak256
	Signature []byte
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
```


### Run

#### Node by node
- Node0
```
go run *.go --config config0.yaml start
```

- Node1
```
go run *.go --config config1.yaml start
```

### Test
```
go test ./...
```

