package zk_circuit

import (
	"github.com/consensys/gnark/frontend"
	"gnark-bid/circuits"
)

type UserData struct {
	UserID      frontend.Variable
	PrivateCode frontend.Variable
}

type Identity struct {
	Nullifier  frontend.Variable `gnark:",public"`
	Commitment frontend.Variable `gnark:",public"`
	Trapdoor   frontend.Variable `gnark:",public"`
}

type BiddingCircuit struct {
	UserData UserData
	Identity Identity `gnark:",public"`

	BidValue frontend.Variable `gnark:",public"`
}

func (circuit *BiddingCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.BidValue), 0)

	// preimage user id
	userHash, err := HashPreImage(api, circuit.UserData.UserID)
	if err != nil {
		return err
	}

	// check user commitment
	userCommitment := circuits.Poseidon(api, []frontend.Variable{userHash, circuit.UserData.PrivateCode})
	api.AssertIsEqual(userCommitment, circuit.Identity.Commitment)

	// check trapdoor
	trapdoor := circuits.Poseidon(api, []frontend.Variable{userCommitment, circuit.Identity.Nullifier})
	api.AssertIsEqual(trapdoor, circuit.Identity.Trapdoor)

	return nil
}
