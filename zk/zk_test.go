package zk_test

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/examples/mimc"
	"github.com/consensys/gnark/test"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"
)

func TestHashPreImage(t *testing.T) {
	assert := test.NewAssert(t)

	preImage := 42
	hash := zk_circuit.HashMIMC(big.NewInt(int64(preImage)).Bytes())
	fmt.Println("hash:", hash.String())
	var circuit mimc.Circuit

	assert.ProverSucceeded(&circuit, &mimc.Circuit{
		Hash:     hash,
		PreImage: preImage,
	}, test.WithCurves(ecc.BN254))
}

func TestMerkleTree(t *testing.T) {
	assert := test.NewAssert(t)

	preImage := 42
	hash := zk_circuit.HashMIMC(big.NewInt(int64(preImage)).Bytes())
	fmt.Println("hash:", hash.String())
	var circuit mimc.Circuit

	assert.ProverSucceeded(&circuit, &mimc.Circuit{
		Hash:     hash,
		PreImage: preImage,
	}, test.WithCurves(ecc.BN254))
}
