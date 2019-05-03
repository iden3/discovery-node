package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// GetSubtopic returns the subtopic string for a giving Address
func GetSubtopic(addr common.Address) string {
	st := fmt.Sprintf("%b", addr.Bytes()[0])
	return st[:4] // 4 is the subtopic addr prefix agroupation
}
