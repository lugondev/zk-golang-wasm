package solidity

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/iden3/go-iden3-crypto/poseidon"
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

	InitSetup[BiddingCircuit]((*ExportSolidityTestSuite[BiddingCircuit])(t), &c, DeployBiddingCircuit, "BiddingCircuit")
}

func (t *ExportSolidityTestSuiteBiddingVerifier) TestVerifyProof() {
	assert := test.NewAssert(t.Suite.T())

	roomId := 1111
	userId := fmt.Sprintf("user=username_%d|room=%d", 2, roomId)
	fmt.Println("user id", userId)
	usernameId := new(big.Int).SetBytes([]byte(userId))

	// create proof bidding
	privateCode := big.NewInt(1111222233334444)
	bidValue := big.NewInt(100) //  keep in mind, don't share this value with anyone
	preHash := zkCircuit.HashMIMC(usernameId.Bytes())
	userCommitment, err := poseidon.Hash([]*big.Int{preHash, privateCode})
	assert.NoError(err, "poseidon hash failed")

	// random nullifier: keep in mind, don't share this value with anyone
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	assert.NoError(err, "generating random bigInt failed")
	nullifier, err := poseidon.HashBytes(buf)
	assert.NoError(err, "generating nullifier failed")
	fmt.Println("nullifier", nullifier.String())

	trapdoorNumber, err := poseidon.Hash([]*big.Int{userCommitment, nullifier})
	assert.NoError(err, "generating trapdoorNumber failed")
	fmt.Println("trapdoor", trapdoorNumber.String())

	merkleAssignment := zkCircuit.BiddingCircuit{
		BidValue: bidValue,
		Identity: zkCircuit.Identity{
			Nullifier:  nullifier,
			Commitment: userCommitment,
			Trapdoor:   trapdoorNumber,
		},
		UserData: zkCircuit.UserData{
			UserID:      usernameId,
			PrivateCode: privateCode,
		},
	}

	var circuit zkCircuit.BiddingCircuit
	assert.ProverSucceeded(&circuit, &merkleAssignment, test.WithCurves(ecc.BN254))

	proofParser, g16Proof, err := t.g16.GenerateProof(&merkleAssignment)
	assert.NoError(err, "proving failed")
	fmt.Println("proof", proofParser)
	fmt.Println("g16Proof", g16Proof)

	// public witness
	var publicInput [4]*big.Int
	publicInput[0] = nullifier
	publicInput[1] = userCommitment
	publicInput[2] = trapdoorNumber
	publicInput[3] = bidValue
	// call the contract
	res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, publicInput)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.True(res, "calling verifier on chain didn't succeed")

	fmt.Println("verifier on chain succeeded")
}

func (t *ExportSolidityTestSuiteBiddingVerifier) TestBidding() {
	assert := test.NewAssert(t.Suite.T())

	bidding, err := zk.NewBidding(nil)
	assert.NoError(err, "creating bidding failed")
	assert.NoError(bidding.InitSession(1111, "username_2", big.NewInt(1111222233334444)), "init session failed")

	// create a proof for a leaf: simulate user 1
	userId := fmt.Sprintf("user=username_%d|room=%d", 2, bidding.RoomID)
	fmt.Println("user id", userId)

	bidValue := big.NewInt(100)
	proof, inputs, err := bidding.GetProof(bidValue)
	assert.NoError(err, "creating proof failed")

	// call the contract
	res, err := t.contract.VerifyProof(nil, proof.A, proof.B, proof.C, inputs)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.True(res, "calling verifier on chain didn't succeed")

	fmt.Println("verifier on chain succeeded")
}
