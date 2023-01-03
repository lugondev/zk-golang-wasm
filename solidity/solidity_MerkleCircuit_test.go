package solidity

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/thoas/go-funk"
	merkle "gnark-bid/merkle"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/consensys/gnark/frontend"
	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteMerkleVerifier ExportSolidityTestSuite[MerkleCircuit]

func TestRunExportSolidityTestSuiteMerkleVerifier(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteMerkleVerifier))
}

func (t *ExportSolidityTestSuiteMerkleVerifier) SetupTest() {
	var c zk_circuit.MerkleCircuit
	c.Path = make([]frontend.Variable, zk.MerkleTreeDepth+1)
	c.Helper = make([]frontend.Variable, zk.MerkleTreeDepth)

	InitSetup[MerkleCircuit]((*ExportSolidityTestSuite[MerkleCircuit])(t), &c, DeployMerkleCircuit, "MerkleCircuit")
}

func (t *ExportSolidityTestSuiteMerkleVerifier) TestVerifyProof() {
	assert := test.NewAssert(t.Suite.T())
	leaves := math.BigPow(2, int64(zk.MerkleTreeDepth))
	var list [][]byte
	for i := 0; i < int(leaves.Int64()); i++ {
		r := fmt.Sprintf("%dabc", i+1)
		list = append(list, []byte(r))
	}

	mkTree, err := merkle.NewMerkleTreeBytesZK(list)
	assert.NoError(err, "creating merkle tree failed")

	proofIndex := 1
	leafHash := mkTree.Hashes[proofIndex]
	fmt.Println("leafHash", hex.EncodeToString(leafHash))
	// create a valid proof
	merkleRoot, merkleProof, proofHelper, err := mkTree.BuilderProofHelper(leafHash)
	assert.NoError(err, "building merkle proof failed")
	fmt.Println("merkleRoot", hex.EncodeToString(merkleRoot))
	fmt.Println("merkleProof length:", len(merkleProof))
	fmt.Println("proofHelper length:", len(proofHelper))
	assert.Equal(len(merkleProof), zk.MerkleTreeDepth+1, "proof length should be equal to zk.MerkleTreeDepth+1")

	verified := merkletree.VerifyProof(mkTree.HashFunc(), merkleRoot, merkleProof, uint64(proofIndex), uint64(mkTree.NumLeaves()))
	assert.True(verified, "merkle proof verification failed")

	merkleAssignment := zk_circuit.MerkleCircuit{
		Path: funk.Map(merkleProof, func(p []byte) frontend.Variable {
			return p
		}).([]frontend.Variable),
		Helper: funk.Map(proofHelper, func(p int) frontend.Variable {
			return p
		}).([]frontend.Variable),
		RootHash: merkleRoot,
	}

	var circuit zk_circuit.MerkleCircuit
	circuit.Path = make([]frontend.Variable, zk.MerkleTreeDepth+1)
	circuit.Helper = make([]frontend.Variable, zk.MerkleTreeDepth)

	assert.ProverSucceeded(&circuit, &merkleAssignment, test.WithCurves(ecc.BN254))

	proofParser, g16Proof, err := t.g16.GenerateProof(&merkleAssignment)
	assert.NoError(err, "proving failed")
	fmt.Println("proof", proofParser)
	fmt.Println("g16Proof", g16Proof)

	// public witness
	var publicInput [1]*big.Int
	publicInput[0] = new(big.Int).SetBytes(merkleRoot)
	// call the contract
	res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, publicInput)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.True(res, "calling verifier on chain didn't succeed")

	// (wrong) public witness
	publicInput[0] = big.NewInt(11)

	// call the contract should fail
	res, err = t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, publicInput)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.False(res, "calling verifier on chain succeed, and shouldn't have")
}
