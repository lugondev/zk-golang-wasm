package zk_circuit

import (
	"github.com/consensys/gnark/frontend"
	"gnark-bid/circuits"
)

type PrivateValueCircuit struct {
	PrivateValue frontend.Variable
	Hash         frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
func (circuit *PrivateValueCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.PrivateValue), 0)

	hashMIMC, err := HashPreImage(api, circuit.PrivateValue)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.Hash, hashMIMC)

	return nil
}
