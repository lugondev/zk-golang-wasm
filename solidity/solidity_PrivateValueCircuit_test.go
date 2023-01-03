package solidity

import (
	"fmt"
	"github.com/consensys/gnark/test"
	"gnark-bid/zk/circuits"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteGroth16 ExportSolidityTestSuite[PrivateValueCircuit]

func TestRunExportSolidityTestSuiteGroth16(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteGroth16))
}

func (t *ExportSolidityTestSuiteGroth16) SetupTest() {
	var c zk_circuit.PrivateValueCircuit
	InitSetup[PrivateValueCircuit]((*ExportSolidityTestSuite[PrivateValueCircuit])(t), &c, DeployPrivateValueCircuit, "PrivateValueCircuit")
}

func (t *ExportSolidityTestSuiteGroth16) TestVerifyProof() {
	assert := test.NewAssert(t.Suite.T())

	pubValue := int64(40)
	privValue := int64(42)
	// create a valid proof
	var assignment zk_circuit.PrivateValueCircuit
	assignment.PrivateValue = privValue
	assignment.Hash = zk_circuit.HashMIMC(big.NewInt(privValue).Bytes())

	proofParser, g16Proof, err := t.g16.GenerateProof(&assignment)
	assert.NoError(err, "proving failed")

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
	assert.NoError(err, "verifying failed")

	// public witness
	input := [1]*big.Int{zk_circuit.HashMIMC(big.NewInt(42).Bytes())}
	// call the contract
	res, err := t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, input)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.True(res, "calling verifier on chain didn't succeed")

	// (wrong) public witness
	input = [1]*big.Int{big.NewInt(pubValue)}

	// call the contract should fail
	res, err = t.contract.VerifyProof(nil, proofParser.A, proofParser.B, proofParser.C, input)
	assert.NoError(err, "calling verifier on chain gave error")
	assert.False(res, "calling verifier on chain succeed, and shouldn't have")
}
