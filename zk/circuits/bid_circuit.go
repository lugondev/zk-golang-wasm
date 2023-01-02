package zk_circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
	"gnark-bid/circuits"
)

type UserData struct {
	PrivateID frontend.Variable
}

type BidCircuit struct {
	UserMerkleRoot                    frontend.Variable `gnark:",public"`
	UserMerkleProof, UserMerkleHelper []frontend.Variable
	User                              UserData

	BidValue frontend.Variable
	BidHash  frontend.Variable `gnark:",public"`
}

func (circuit *BidCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.BidValue), 0)

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	merkle.VerifyProof(api, hFunc, circuit.UserMerkleRoot, circuit.UserMerkleProof, circuit.UserMerkleHelper)

	leaf, err := HashPreImage(api, circuit.User.PrivateID)
	if err != nil {
		return err
	}
	api.AssertIsEqual(leaf, circuit.UserMerkleProof[0])

	hashMIMC, err := HashPreImage(api, circuit.BidValue)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.BidHash, hashMIMC)
	return nil
}
