package solidity

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	merkle "gnark-bid/merkle"
	"gnark-bid/wasm"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
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
	vk      groth16.VerifyingKey
	pk      groth16.ProvingKey
	proof   *bytes.Buffer
	circuit zk_circuit.MerkleCircuit
	r1cs    frontend.CompiledConstraintSystem
}

func TestRunExportSolidityTestSuiteMerkleVerifier(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteMerkleVerifier))
}

var lenProof = 10

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

	t.circuit = zk_circuit.MerkleCircuit{
		Path:   make([]frontend.Variable, lenProof),
		Helper: make([]frontend.Variable, lenProof-1),
	}
	t.r1cs, err = frontend.Compile(ecc.BN254, r1cs.NewBuilder, &t.circuit)
	t.NoError(err, "compiling R1CS failed")

	vpKey, err := wasm.GetVPKey("MerkleCircuit")
	t.NoError(err, "getting vpkey failed")
	// read proving and verifying keys
	t.pk = vpKey.PK
	t.vk = vpKey.VK

}

func (t *ExportSolidityTestSuiteMerkleVerifier) TestVerifyProof() {
	leaves := math.BigPow(2, int64(lenProof))
	var list [][]byte
	for i := 0; i < int(leaves.Int64()); i++ {
		r := fmt.Sprintf("%dabc", i+1)
		list = append(list, []byte(r))
	}
	mkTree, err := merkle.NewMerkleTreeBytes(list)
	t.NoError(err, "creating merkle tree failed")

	proofIndex := 6
	leafHash := mkTree.Hashes[proofIndex]
	fmt.Println("leafHash", hex.EncodeToString(leafHash))
	merkleRoot, merkleProof, proofHelper, err := mkTree.BuilderProofHelper(leafHash)
	t.NoError(err, "building merkle proof failed")
	// create a valid proof

	merkleAssignment := &zk_circuit.MerkleCircuit{
		Path:     make([]frontend.Variable, lenProof),
		Helper:   make([]frontend.Variable, lenProof-1),
		RootHash: merkleRoot,
	}
	for i := 0; i < lenProof; i++ {
		merkleAssignment.Path[i] = merkleProof[i]
	}
	for i := 0; i < lenProof-1; i++ {
		merkleAssignment.Helper[i] = proofHelper[i]
	}

	// witness creation
	witness, err := frontend.NewWitness(merkleAssignment, ecc.BN254)
	t.NoError(err, "witness creation failed")

	// prove
	proof, err := groth16.Prove(t.r1cs, t.pk, witness)
	t.NoError(err, "proving failed")
	{
		_, err = proof.WriteRawTo(t.proof)
		t.NoError(err, "writing proof failed")
	}

	// witness creation
	hiddenWitness, err := frontend.NewWitness(merkleAssignment, ecc.BN254)
	// ensure gnark (Go) code verifies it
	publicWitness, _ := hiddenWitness.Public()
	fmt.Println("publicWitness:", publicWitness)

	err = groth16.Verify(proof, t.vk, publicWitness)
	t.NoError(err, "verifying failed")

	//var buf bytes.Buffer
	//_, _ = proof.WriteRawTo(&buf)
	//proofBytes := buf.Bytes()
	//
	//proofParser := zk.ParserProof(proofBytes)
	//
	//// public witness
	//proofParser.Input[0] = zk_circuit.HashMIMC(big.NewInt(42).Bytes())
	//// call the contract
	//res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, proofParser.Input)
	//if t.NoError(err, "calling verifier on chain gave error") {
	//	t.True(res, "calling verifier on chain didn't succeed")
	//}
	//
	//// (wrong) public witness
	//proofParser.Input[0] = big.NewInt(pubValue)
	//
	//// call the contract should fail
	//res, err = t.verifierContract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, proofParser.Input)
	//if t.NoError(err, "calling verifier on chain gave error") {
	//	t.False(res, "calling verifier on chain succeed, and shouldn't have")
	//}
}
