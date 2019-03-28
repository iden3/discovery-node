package discovery

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// DiscoveryService is a concretion of Service data type, for the discovery-node service
type DiscoveryService Service

// NewDiscoveryService creates a new DiscoveryService
func NewDiscoveryService(idAddr common.Address, kademliaAddr []byte, pssPubK *ecdsa.PublicKey, url, mode string, proofService []byte) (DiscoveryService, error) {
	d := DiscoveryService{
		IdAddr:       idAddr,
		KademliaAddr: kademliaAddr,
		PssPubK:      PubK{*pssPubK},
		Url:          url,
		Type:         DISCOVERYTYPE,
		Mode:         mode, // Active or Passive
		ProofService: proofService,
	}

	return d, nil
}

// DiscoverIdentity generates the Query about an identity and sends it over Swarm Pss
func (d *DiscoveryService) NewQueryPacket(idAddr common.Address) (*Query, error) {
	q := &Query{
		Version:          DISCOVERYVERSION,
		AboutId:          idAddr,
		RequesterId:      d.IdAddr,
		RequesterKAddr:   d.KademliaAddr,
		RequesterPssPubK: d.PssPubK,
		InfoFrom:         []byte{},
		Timestamp:        time.Now().Unix(),
		Nonce:            0,
	}

	// TODO calculate PoW

	// TODO sign packet

	return q, nil
}

// AnswerRequest generates and returns the answer for a Query request for which knows the answer
// first, the Discovery Node will check if knows the answer
func (d *DiscoveryService) NewAnswerPacket(q *Query, id *Id) (*Answer, error) {
	// check that the query and id are about the same idaddr
	if !bytes.Equal(q.AboutId.Bytes(), id.IdAddr.Bytes()) {
		return nil, errors.New("resolved idAddr is not the same than query.IdAddr")
	}

	// generate the answer data packet
	answer := &Answer{
		Version:   DISCOVERYVERSION,
		AboutId:   q.AboutId,
		FromId:    d.IdAddr,
		AgentId:   Service(*d),
		Services:  id.Services, // TODO data related to the requested idAddr
		Timestamp: time.Now().Unix(),
		Signature: []byte{},
	}

	// TODO sign packet

	return answer, nil
}
