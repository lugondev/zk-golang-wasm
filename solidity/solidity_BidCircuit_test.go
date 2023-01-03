package solidity

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/thoas/go-funk"
	merkle "gnark-bid/merkle"
	"gnark-bid/zk"
	zkCircuit "gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteBiddingVerifier ExportSolidityTestSuite[BiddingCircuit]

func TestRunExportSolidityTestSuiteBiddingVerifier(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteBiddingVerifier))
}

func (t *ExportSolidityTestSuiteBiddingVerifier) SetupTest() {

	var c zkCircuit.BiddingCircuit
	c.UserMerklePath = make([]frontend.Variable, zk.MerkleTreeDepth+1)
	c.UserMerkleHelper = make([]frontend.Variable, zk.MerkleTreeDepth)

	InitSetup[BiddingCircuit]((*ExportSolidityTestSuite[BiddingCircuit])(t), &c, DeployBiddingCircuit, "BiddingCircuit")
}

func (t *ExportSolidityTestSuiteBiddingVerifier) TestVerifyProof() {
	assert := test.NewAssert(t.Suite.T())

	leaves := math.BigPow(2, int64(zk.MerkleTreeDepth))
	var list [][]byte
	// generate leaves: users use this to generate their own merkle tree
	for i := 0; i < int(leaves.Int64()); i++ {
		r := fmt.Sprintf("username_%d", i+1)
		list = append(list, []byte(r))
	}

	// create ZK Tree with MIMC hash function
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

	// verify proof
	verified := merkletree.VerifyProof(mkTree.HashFunc(), merkleRoot, merkleProof, uint64(proofIndex), uint64(mkTree.NumLeaves()))
	assert.True(verified, "merkle proof verification failed")

	merkleAssignment := zkCircuit.BiddingCircuit{
		UserMerklePath: funk.Map(merkleProof, func(p []byte) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleHelper: funk.Map(proofHelper, func(p int) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleRoot: merkleRoot,
	}

	var circuit zkCircuit.BiddingCircuit
	circuit.UserMerklePath = make([]frontend.Variable, zk.MerkleTreeDepth+1)
	circuit.UserMerkleHelper = make([]frontend.Variable, zk.MerkleTreeDepth)

	assert.ProverSucceeded(&circuit, &merkleAssignment, test.WithCurves(ecc.BN254))

	proofParser, g16Proof, err := t.g16.GenerateProof(&merkleAssignment)
	assert.NoError(err, "proving failed")
	fmt.Println("proof", proofParser)
	fmt.Println("g16Proof", g16Proof)

	// public witness
	var publicInput [2]*big.Int
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
