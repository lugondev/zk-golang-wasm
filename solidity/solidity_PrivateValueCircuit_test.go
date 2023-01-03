package solidity

import (
	"fmt"
	"gnark-bid/wasm"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteGroth16 struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend

	// verifier contract
	contract *PrivateValueCircuit

	// groth16 gnark objects
	vk groth16.VerifyingKey
	pk groth16.ProvingKey

	g16 *zk.GnarkGroth16
}

func TestRunExportSolidityTestSuiteGroth16(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteGroth16))
}

func (t *ExportSolidityTestSuiteGroth16) SetupTest() {

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
	_, _, v, err := DeployPrivateValueCircuit(auth, t.backend)
	t.NoError(err, "deploy verifier contract failed")
	t.contract = v
	t.backend.Commit()

	vpKey, err := wasm.GetVPKey("PrivateValueCircuit")
	t.NoError(err, "getting vpkey failed")

	// read proving and verifying keys
	t.pk = vpKey.PK
	t.vk = vpKey.VK

	var c zk_circuit.PrivateValueCircuit
	t.g16, err = zk.NewGnarkGroth16(vpKey, &c)
	t.NoError(err, "init gnark groth16 failed")
}

func (t *ExportSolidityTestSuiteGroth16) TestVerifyProof() {

	pubValue := int64(40)
	privValue := int64(42)
	// create a valid proof
	var assignment zk_circuit.PrivateValueCircuit
	assignment.PrivateValue = privValue
	assignment.Hash = zk_circuit.HashMIMC(big.NewInt(privValue).Bytes())

	proofParser, g16Proof, err := t.g16.GenerateProof(&assignment)
	t.NoError(err, "proving failed")

	// hidden witness
	var hiddenAssignment zk_circuit.PrivateValueCircuit
	hiddenAssignment.PrivateValue = int64(0)
	hiddenAssignment.Hash = zk_circuit.HashMIMC(big.NewInt(privValue).Bytes())

	// witness creation
	hiddenWitness, err := frontend.NewWitness(&hiddenAssignment, ecc.BN254)
	// ensure gnark (Go) code verifies it
	publicWitness, _ := hiddenWitness.Public()
	fmt.Println("publicWitness:", publicWitness)

	err = groth16.Verify(g16Proof, t.vk, publicWitness)
	t.NoError(err, "verifying failed")

	// public witness
	input := [1]*big.Int{zk_circuit.HashMIMC(big.NewInt(42).Bytes())}
	// call the contract
	res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.True(res, "calling verifier on chain didn't succeed")
	}

	// (wrong) public witness
	input = [1]*big.Int{big.NewInt(pubValue)}

	// call the contract should fail
	res, err = t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.False(res, "calling verifier on chain succeed, and shouldn't have")
	}
}
