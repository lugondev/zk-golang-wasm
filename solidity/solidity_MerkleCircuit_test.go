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
	"gnark-bid/wasm"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteMerkleVerifier struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend

	// contract
	contract *MerkleCircuit

	// groth16 gnark objects
	vk  groth16.VerifyingKey
	pk  groth16.ProvingKey
	g16 *zk.GnarkGroth16
}

func TestRunExportSolidityTestSuiteMerkleVerifier(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteMerkleVerifier))
}

var lenProof = zk.MerkleTreeDepth

func (t *ExportSolidityTestSuiteMerkleVerifier) SetupTest() {

	const gasLimit uint64 = 4712388

	// setup simulated backend
	key, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	t.NoError(err, "init keyed transactor")

	genesis := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(1000000000000000000)}, // 1 Eth
	}
	t.backend = backends.NewSimulatedBackend(genesis, gasLimit)

	// deploy verifier contract
	_, _, v, err := DeployMerkleCircuit(auth, t.backend)
	t.NoError(err, "deploy verifier contract failed")
	t.contract = v
	t.backend.Commit()

	vpKey, err := wasm.GetVPKey("MerkleCircuit")
	t.NoError(err, "getting vpkey failed")

	// read proving and verifying keys
	t.pk = vpKey.PK
	t.vk = vpKey.VK

	var c zk_circuit.MerkleCircuit
	c.Path = make([]frontend.Variable, lenProof+1)
	c.Helper = make([]frontend.Variable, lenProof)

	t.g16, err = zk.NewGnarkGroth16(vpKey, &c)
	t.NoError(err, "init groth16 failed")
}

func (t *ExportSolidityTestSuiteMerkleVerifier) TestVerifyProof() {
	assert := test.NewAssert(t.T())
	leaves := math.BigPow(2, int64(lenProof))
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
	assert.Equal(len(merkleProof), lenProof+1, "proof length should be equal to lenProof+1")

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
	circuit.Path = make([]frontend.Variable, lenProof+1)
	circuit.Helper = make([]frontend.Variable, lenProof)

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
	if t.NoError(err, "calling verifier on chain gave error") {
		t.True(res, "calling verifier on chain didn't succeed")
	}

	// (wrong) public witness
	publicInput[0] = big.NewInt(11)

	// call the contract should fail
	res, err = t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, publicInput)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.False(res, "calling verifier on chain succeed, and shouldn't have")
	}
}
