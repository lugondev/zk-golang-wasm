package zk_circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
	"gnark-bid/circuits"
)

type UserData struct {
	UserID    frontend.Variable
	PrivateID frontend.Variable
}

type DataID struct {
	Nullifier    frontend.Variable `gnark:",public"`
	IdCommitment frontend.Variable `gnark:",public"`
	Trapdoor     frontend.Variable `gnark:",public"`
}

type BiddingCircuit struct {
	UserMerkleRoot                   frontend.Variable `gnark:",public"`
	UserMerklePath, UserMerkleHelper []frontend.Variable

	UserData UserData
	DataID   DataID `gnark:",public"`

	BidValue frontend.Variable `gnark:",public"`
}

func (circuit *BiddingCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.BidValue), 0)

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	merkle.VerifyProof(api, hFunc, circuit.UserMerkleRoot, circuit.UserMerklePath, circuit.UserMerkleHelper)

	// preimage user id
	userHash, err := HashPreImage(api, circuit.UserData.UserID)
	if err != nil {
		return err
	}

	api.Println(circuit.UserMerklePath[0], userHash)
	// check user in merkle tree
	api.AssertIsEqual(circuit.UserMerklePath[0], userHash)

	// check user commitment
	userCommitment := circuits.Poseidon(api, []frontend.Variable{circuit.UserData.UserID, circuit.UserData.PrivateID})
	api.AssertIsEqual(userCommitment, circuit.DataID.IdCommitment)

	// check trapdoor
	trapdoor := circuits.Poseidon(api, []frontend.Variable{userCommitment, circuit.DataID.Nullifier})
	api.AssertIsEqual(trapdoor, circuit.DataID.Trapdoor)

	return nil
}
