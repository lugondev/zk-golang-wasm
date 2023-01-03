package zk

import (
	"github.com/consensys/gnark/backend/groth16"
	"math/big"
)

type Proof struct {
	A [2]*big.Int    `json:"a"`
	B [2][2]*big.Int `json:"b"`
	C [2]*big.Int    `json:"c"`
}

type PublicInput []*big.Int

type VPKey struct {
	ProvingKey   string `json:"provingKey"`
	VerifyingKey string `json:"verifyingKey"`
	VK           groth16.VerifyingKey
	PK           groth16.ProvingKey
}
