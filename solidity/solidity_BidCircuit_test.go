package solidity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/iden3/go-iden3-crypto/poseidon"
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
		list = append(list, new(big.Int).SetBytes([]byte(r)).Bytes())
	}

	// create ZK Tree with MIMC hash function
	mkTree, err := merkle.NewMerkleTreeBytesZK(list)
	assert.NoError(err, "creating merkle tree failed")

	// create a proof for a leaf: simulate user 1
	username := fmt.Sprintf("username_%d", 2)
	fmt.Println("username", username)

	usernameId := new(big.Int).SetBytes([]byte(username))
	leafHash := zkCircuit.HashMIMC(usernameId.Bytes()).Bytes()
	fmt.Println("hash", hex.EncodeToString(leafHash))
	proofIndex := funk.IndexOf(mkTree.Hashes, leafHash)
	fmt.Println("proofIndex", proofIndex)

	// create a valid proof
	merkleRoot, merkleProof, proofHelper, err := mkTree.BuilderProofHelper(leafHash)
	assert.NoError(err, "building merkle proof failed")
	assert.Equal(len(merkleProof), zk.MerkleTreeDepth+1, "proof length should be equal to zk.MerkleTreeDepth+1")
	for i := range merkleProof {
		fmt.Println("merkleProof", i, hex.EncodeToString(merkleProof[i]))
	}

	// verify proof: valid user
	verified := merkletree.VerifyProof(mkTree.HashFunc(), merkleRoot, merkleProof, uint64(proofIndex), uint64(mkTree.NumLeaves()))
	assert.True(verified, "merkle proof verification failed")

	// create proof bidding
	privateUserId := big.NewInt(1111222233334444)
	bidValue := big.NewInt(100) //  keep in mind, don't share this value with anyone
	idCommitment, err := poseidon.Hash([]*big.Int{usernameId, privateUserId})
	assert.NoError(err, "poseidon hash failed")

	// random nullifier: keep in mind, don't share this value with anyone
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	assert.NoError(err, "generating random bigInt failed")
	nullifier, err := poseidon.HashBytes(buf)
	assert.NoError(err, "generating nullifier failed")
	fmt.Println("nullifier", nullifier.String())

	trapdoorNumber, err := poseidon.Hash([]*big.Int{idCommitment, nullifier})
	assert.NoError(err, "generating trapdoorNumber failed")
	fmt.Println("trapdoor", trapdoorNumber.String())

	merkleAssignment := zkCircuit.BiddingCircuit{
		UserMerklePath: funk.Map(merkleProof, func(p []byte) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleHelper: funk.Map(proofHelper, func(p int) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleRoot: merkleRoot,
		BidValue:       bidValue,
		DataID: zkCircuit.DataID{
			Nullifier:    nullifier,
			IdCommitment: idCommitment,
			Trapdoor:     trapdoorNumber,
		},
		UserData: zkCircuit.UserData{
			UserID:    usernameId,
			PrivateID: privateUserId,
		},
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
	var publicInput [5]*big.Int
	publicInput[0] = new(big.Int).SetBytes(merkleRoot)
	publicInput[1] = nullifier
	publicInput[2] = idCommitment
	publicInput[3] = trapdoorNumber
	publicInput[4] = bidValue
	// call the contract
	res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, publicInput)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.True(res, "calling verifier on chain didn't succeed")

	fmt.Println("verifier on chain succeeded")
}
