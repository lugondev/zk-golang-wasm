package zk

import (
	"github.com/consensys/gnark/frontend"
	"gnark-bid/circuits"
)

type BidCircuit struct {
	PrivateValue frontend.Variable
	Hash         frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
func (circuit *BidCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.PrivateValue), 0)

	hashMIMC := HashPreImage(api, circuit.PrivateValue)
	api.AssertIsEqual(circuit.Hash, hashMIMC)

	return nil
}
