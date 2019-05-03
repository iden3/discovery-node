package node

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const addrA0Hex = "0xaa0cf716f7f176f29DFbfc792D7B0Bf097455777"
const addrA1Hex = "0xaa166b8dc6762A972eC4129abcA313C9F19861E6"
const addrA2Hex = "0xaa271E0D634c64E85a1362599f11eA02A781Eb8C"
const addrB0Hex = "0xbb06D140E7e5F0cF9799765A5836bA523278cdC3"
const addrB1Hex = "0xbb1eFf3144201D7c57e12937CF5c0D5715221BF3"
const addrC0Hex = "0xcc09c5d2603B2C7bdEFbeB9735cc3De6C58b20a2"

func TestSubtopic(t *testing.T) {
	stA0 := GetSubtopic(common.HexToAddress(addrA0Hex))
	stA1 := GetSubtopic(common.HexToAddress(addrA1Hex))
	stA2 := GetSubtopic(common.HexToAddress(addrA2Hex))
	stB0 := GetSubtopic(common.HexToAddress(addrB0Hex))
	stB1 := GetSubtopic(common.HexToAddress(addrB1Hex))
	stC0 := GetSubtopic(common.HexToAddress(addrC0Hex))

	assert.Equal(t, stA0, "1010")
	assert.Equal(t, stA0, stA1)
	assert.Equal(t, stA0, stA2)
	assert.Equal(t, stB0, stB1)
	assert.NotEqual(t, stA0, stB0)
	assert.NotEqual(t, stB0, stC0)
	assert.NotEqual(t, stA0, stC0)
}
