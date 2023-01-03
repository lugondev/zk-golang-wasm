package zk

import (
	"bytes"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

type GnarkGroth16 struct {
	vk   groth16.VerifyingKey
	pk   groth16.ProvingKey
	r1cs frontend.CompiledConstraintSystem
}

func NewGnarkGroth16(key *VPKey, circuit frontend.Circuit) (*GnarkGroth16, error) {
	g16 := &GnarkGroth16{}

	if _r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit); err != nil {
		log.Println(err, "compiling R1CS failed")
		return nil, err
	} else {
		g16.r1cs = _r1cs
	}

	if err := g16.setup(key); err != nil {
		return nil, err
	}
	return g16, nil
}

func (t *GnarkGroth16) setup(vpKey *VPKey) error {
	// read proving and verifying keys
	t.pk = groth16.NewProvingKey(ecc.BN254)
	{
		pkBuf := bytes.NewBuffer(common.FromHex(vpKey.ProvingKey))
		if _, err := t.pk.ReadFrom(pkBuf); err != nil {
			log.Println(err, "reading proving key failed")
			return err
		}
	}
	t.vk = groth16.NewVerifyingKey(ecc.BN254)
	{
		vkBuf := bytes.NewBuffer(common.FromHex(vpKey.VerifyingKey))
		if _, err := t.vk.ReadFrom(vkBuf); err != nil {
			log.Println(err, "reading verifying key failed")
			return err
		}
	}
	return nil
}

func (t *GnarkGroth16) GenerateProof(assignment frontend.Circuit) (*Proof, groth16.Proof, error) {
	// witness creation
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		return nil, nil, err
	}

	// prove
	proof, err := groth16.Prove(t.r1cs, t.pk, witness)
	if err != nil {

		return nil, nil, err
	}

	// ensure gnark (Go) code verifies it
	publicWitness, err := witness.Public()
	if err != nil {
		return nil, nil, err
	}

	if err := groth16.Verify(proof, t.vk, publicWitness); err != nil {
		return nil, nil, err
	}

	// get proof bytes
	var proofBuffer bytes.Buffer
	if _, err := proof.WriteRawTo(&proofBuffer); err != nil {
		return nil, nil, err
	}
	proofBytes := proofBuffer.Bytes()
	proofStruct := ParserProof(proofBytes)

	return proofStruct, proof, nil
}
