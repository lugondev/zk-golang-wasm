package zk_test

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/examples/mimc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/thoas/go-funk"
	"gnark-bid/merkle"
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

	var circuit zk_circuit.MerkleCircuit
	tree, err := merkle_tree.NewMerkleTreeFromBytes([][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	})
	assert.NoError(err)

	proofs, err := tree.GetProofFromByte([]byte("a"))
	assert.NoError(err)
	path := funk.Map(proofs, func(p []byte) frontend.Variable {
		return p
	}).([]frontend.Variable)

	leaf, err := tree.GetLeafFromByte([]byte("a"))
	assert.NoError(err)
	assert.ProverSucceeded(&circuit, &zk_circuit.MerkleCircuit{
		RootHash: tree.GetRoot(),
		Path:     path,
		Helper:   []frontend.Variable{leaf},
	}, test.WithCurves(ecc.BN254))
}
