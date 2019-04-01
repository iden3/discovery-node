package node

import "github.com/iden3/discovery-node/utils"

func (node *NodeSrv) SignBytes(b []byte) ([]byte, error) {
	h := utils.EthHash(b)
	sig, err := node.ks.SignHash(node.acc, h)
	if err != nil {
		return []byte{}, err
	}
	sig[64] += 27
	return sig, nil
}