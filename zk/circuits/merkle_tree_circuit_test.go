package zk_circuit_test

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	gnarkMerkle "github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/test"
	"github.com/ethereum/go-ethereum/common"
	merkle "gnark-bid/merkle"
	"testing"
)

type merkleCircuit struct {
	RootHash     frontend.Variable `gnark:",public"`
	Path, Helper []frontend.Variable
}

func (circuit *merkleCircuit) Define(api frontend.API) error {
	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	gnarkMerkle.VerifyProof(api, hFunc, circuit.RootHash, circuit.Path, circuit.Helper)
	return nil
}

func TestMerkleCircuit(t *testing.T) {
	assert := test.NewAssert(t)

	var list [][]byte
	for i := 0; i < 1024*2; i++ {
		r := fmt.Sprintf("%dabc", i+1)
		list = append(list, []byte(r))
	}

	mkTree, err := merkle.NewMerkleTreeBytes(list)
	assert.NoError(err)

	proofIndex := 6
	leafHash := mkTree.Hashes[proofIndex]
	fmt.Println("leafHash", hex.EncodeToString(leafHash))
	merkleRoot, merkleProof, proofHelper, err := mkTree.BuilderProofHelper(leafHash)
	assert.NoError(err)

	fmt.Println("merkleRoot:", common.Bytes2Hex(merkleRoot))
	for i, p := range merkleProof {
		fmt.Println("merkleProof:", i, common.Bytes2Hex(p))
	}
	fmt.Println("proofHelper:", proofHelper)

	verified := merkletree.VerifyProof(hash.MIMC_BN254.New(), merkleRoot, merkleProof, uint64(proofIndex), uint64(mkTree.NumLeaves()))
	if !verified {
		fmt.Printf("The merkle proof in plain go should pass")
	}
	assert.True(verified)

	// create cs
	circuit := merkleCircuit{
		Path:   make([]frontend.Variable, len(merkleProof)),
		Helper: make([]frontend.Variable, len(merkleProof)-1),
	}
	_r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	assert.NoError(err)

	pk, vk, err := groth16.Setup(_r1cs)
	assert.NoError(err)

	merkleCc := &merkleCircuit{
		Path:     make([]frontend.Variable, len(merkleProof)),
		Helper:   make([]frontend.Variable, len(merkleProof)-1),
		RootHash: merkleRoot,
	}
	for i := 0; i < len(merkleProof); i++ {
		merkleCc.Path[i] = merkleProof[i]
	}
	for i := 0; i < len(merkleProof)-1; i++ {
		merkleCc.Helper[i] = proofHelper[i]
	}

	witness, err := frontend.NewWitness(merkleCc, ecc.BN254)
	assert.NoError(err, "witness creation failed")

	proof, err := groth16.Prove(_r1cs, pk, witness)
	assert.NoError(err)

	publicWitness, _ := witness.Public()
	err = groth16.Verify(proof, vk, publicWitness)
	assert.NoError(err)
	fmt.Printf("verification succeded\n")
}
