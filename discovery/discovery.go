package discovery

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type DiscoveryService Service

func NewDiscoveryService(idAddr common.Address, pubK *ecdsa.PublicKey, url, mode string, proofServer []byte) (DiscoveryService, error) {
	d := DiscoveryService{
		IdAddr:      idAddr,
		PubK:        pubK,
		Url:         url,
		Type:        DISCOVERYTYPE,
		Mode:        mode, // Active or Passive
		ProofServer: proofServer,
	}

	return d, nil
}

// DiscoverIdentity generates the Query about an identity and sends it over Swarm Pss
func (d *DiscoveryService) NewQueryPacket(idAddr common.Address) (Query, error) {
	q := Query{
		About:     idAddr,
		From:      d.IdAddr,
		InfoFrom:  []byte{},
		Nonce:     0,
		PoW:       [32]byte{}, // TODO
		Signature: []byte{},   // TODO
	}
	fmt.Println(q)

	// TODO calculate PoW

	// TODO sign packet

	return q, nil
}

// AnswerRequest generates and returns the answer for a Query request for which knows the answer
// first, the Discovery Node will check if knows the answer
func (d *DiscoveryService) NewAnswerPacket(q Query) {
	// get data from the requested q.About

	// generate the answer data packet
	answer := Answer{
		About:     q.About,
		From:      d.IdAddr,
		AgentId:   Service(*d),
		Services:  []Service{}, // TODO data related to the requested idAddr
		Signature: []byte{},
	}
	fmt.Println(answer)

	// TODO sign packet

	// send Answer packet

}
