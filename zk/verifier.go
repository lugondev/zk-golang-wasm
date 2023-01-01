package zk

import (
	"bytes"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/common"
	"gnark-bid/zk/circuits"
	"log"
	"math/big"
)

type GnarkGroth16 struct {
	vk   groth16.VerifyingKey
	pk   groth16.ProvingKey
	r1cs frontend.CompiledConstraintSystem
}

func NewGnarkGroth16(key *VPKey, circuit frontend.Circuit) (*GnarkGroth16, error) {
	g16 := &GnarkGroth16{}
	err := g16.setup(key, circuit)
	if err != nil {
		return nil, err
	}
	return g16, nil
}

func (t *GnarkGroth16) setup(vpKey *VPKey, circuit frontend.Circuit) error {
	var err error
	t.r1cs, err = frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		log.Println(err, "compiling R1CS failed")
		return err
	}
	if err != nil {
		log.Println(err, "reading file failed")
		return err
	}
	// read proving and verifying keys
	t.pk = groth16.NewProvingKey(ecc.BN254)
	{
		pkBuf := bytes.NewBuffer(common.FromHex(vpKey.ProvingKey))
		_, err = t.pk.ReadFrom(pkBuf)
		if err != nil {
			log.Println(err, "reading proving key failed")
			return err
		}
	}
	t.vk = groth16.NewVerifyingKey(ecc.BN254)
	{
		vkBuf := bytes.NewBuffer(common.FromHex(vpKey.VerifyingKey))
		_, err = t.vk.ReadFrom(vkBuf)
		if err != nil {
			log.Println(err, "reading verifying key failed")
			return err
		}
	}
	return nil
}

func (t *GnarkGroth16) GenerateProof(assignment zk_circuit.PrivateValueCircuit, input [1]*big.Int) (*Proof, error) {
	// witness creation
	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		return nil, err
	}

	// prove
	proof, err := groth16.Prove(t.r1cs, t.pk, witness)
	if err != nil {
		return nil, err
	}

	// ensure gnark (Go) code verifies it
	publicWitness, _ := witness.Public()

	err = groth16.Verify(proof, t.vk, publicWitness)
	if err != nil {
		return nil, err
	}

	// get proof bytes
	var buf bytes.Buffer
	_, err = proof.WriteRawTo(&buf)
	if err != nil {
		return nil, err
	}
	proofBytes := buf.Bytes()
	proofStruct := ParserProof(proofBytes)
	proofStruct.Input = input

	return proofStruct, nil
}
