package solidity

import (
	"bytes"
	"fmt"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"math/big"
	"os"
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

type ExportSolidityTestSuiteGroth16 struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend

	// verifier contract
	verifierContract *Verifier

	// groth16 gnark objects
	vk      groth16.VerifyingKey
	pk      groth16.ProvingKey
	proof   *bytes.Buffer
	circuit zk_circuit.PrivateValueCircuit
	r1cs    frontend.CompiledConstraintSystem
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
	_, _, v, err := DeployVerifier(auth, t.backend)
	t.NoError(err, "deploy verifier contract failed")
	t.verifierContract = v
	t.backend.Commit()

	t.r1cs, err = frontend.Compile(ecc.BN254, r1cs.NewBuilder, &t.circuit)
	t.NoError(err, "compiling R1CS failed")

	// read proving and verifying keys
	t.pk = groth16.NewProvingKey(ecc.BN254)
	{
		f, _ := os.Open("zk.g16.pk")
		_, err = t.pk.ReadFrom(f)
		_ = f.Close()
		t.NoError(err, "reading proving key failed")
	}

	t.vk = groth16.NewVerifyingKey(ecc.BN254)
	{
		f, _ := os.Open("zk.g16.vk")
		_, err = t.vk.ReadFrom(f)
		_ = f.Close()
		t.NoError(err, "reading verifying key failed")
	}

}

func (t *ExportSolidityTestSuiteGroth16) TestVerifyProof() {

	pubValue := int64(40)
	privValue := int64(42)
	// create a valid proof
	var assignment zk_circuit.PrivateValueCircuit
	assignment.PrivateValue = privValue
	assignment.Hash = zk_circuit.HashMIMC(big.NewInt(privValue).Bytes())

	// witness creation
	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	t.NoError(err, "witness creation failed")

	// prove
	proof, err := groth16.Prove(t.r1cs, t.pk, witness)
	t.NoError(err, "proving failed")
	{
		_, err = proof.WriteRawTo(t.proof)
		t.NoError(err, "writing proof failed")
	}

	// hidden witness
	var hiddenAssignment zk_circuit.PrivateValueCircuit
	hiddenAssignment.PrivateValue = int64(0)
	hiddenAssignment.Hash = zk_circuit.HashMIMC(big.NewInt(privValue).Bytes())

	// witness creation
	hiddenWitness, err := frontend.NewWitness(&hiddenAssignment, ecc.BN254)
	// ensure gnark (Go) code verifies it
	publicWitness, _ := hiddenWitness.Public()
	fmt.Println("publicWitness:", publicWitness)

	err = groth16.Verify(proof, t.vk, publicWitness)
	t.NoError(err, "verifying failed")

	var buf bytes.Buffer
	_, _ = proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	proofParser := zk.ParserProof(proofBytes)

	// public witness
	proofParser.Input[0] = zk_circuit.HashMIMC(big.NewInt(42).Bytes())
	// call the contract
	res, err := t.verifierContract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, proofParser.Input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.True(res, "calling verifier on chain didn't succeed")
	}

	// (wrong) public witness
	proofParser.Input[0] = big.NewInt(pubValue)

	// call the contract should fail
	res, err = t.verifierContract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, proofParser.Input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.False(res, "calling verifier on chain succeed, and shouldn't have")
	}
}
