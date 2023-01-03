package zk_test

import (
	"encoding/hex"
	"fmt"
	gnarkMerkleTree "github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/examples/mimc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/thoas/go-funk"
	merkle "gnark-bid/merkle"
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

	tree, err := merkle.NewMerkleTreeBytesZK([][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	})
	assert.NoError(err)
	for _, h := range tree.Hashes {
		fmt.Println("hash:", hex.EncodeToString(h))
	}
	leafHash, err := tree.GetLeaf(merkle.ByteContent{B: []byte("b")})
	assert.NoError(err)
	fmt.Println("root:", tree.RootHex())
	fmt.Println("leafHash:", hex.EncodeToString(leafHash))
	merkleRoot, merkleProof, proofHelper, err := tree.BuilderProofHelper(leafHash)
	fmt.Println("merkleRoot:", hex.EncodeToString(merkleRoot))
	assert.NoError(err)
	fmt.Println("proof length:", len(merkleProof))
	fmt.Println("helper length:", len(proofHelper))
	assert.True(gnarkMerkleTree.VerifyProof(tree.HashFunc(), merkleRoot, merkleProof, 1, 4), "proof should be valid")

	var circuit zk_circuit.MerkleCircuit
	circuit.Path = make([]frontend.Variable, 3)
	circuit.Helper = make([]frontend.Variable, 2)
	proof := funk.Map(merkleProof, func(p []byte) frontend.Variable {
		return p
	}).([]frontend.Variable)

	merkleAssignment := &zk_circuit.MerkleCircuit{
		RootHash: merkleRoot,
		Path:     proof,
		Helper: funk.Map(proofHelper, func(p int) frontend.Variable {
			return p
		}).([]frontend.Variable),
	}

	assert.ProverSucceeded(&circuit, merkleAssignment, test.WithCurves(ecc.BN254))

	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	assert.NoError(err, "compilation failed")
	pk, vk, err := groth16.Setup(_r1cs)
	assert.NoError(err, "setup failed")
	witness, err := frontend.NewWitness(merkleAssignment, ecc.BN254)
	assert.NoError(err, "failed to create witness")

	fmt.Println("GenerateProof 2:", "prove process")
	// prove
	g16Proof, err := groth16.Prove(_r1cs, pk, witness)
	assert.NoError(err, "failed to prove")
	fmt.Println("g16Proof:", g16Proof)
	publicWitness, err := witness.Public()
	assert.NoError(err, "failed to get public witness")

	if err := groth16.Verify(g16Proof, vk, publicWitness); err != nil {
		assert.NoError(err, "failed to verify")
	} else {
		fmt.Println("Verify success")
	}
}
