package zk

import (
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"math/big"
)

func HashMIMC(pre []byte) *big.Int {
	h := hash.MIMC_BN254.New()
	h.Write(pre)

	return new(big.Int).SetBytes(h.Sum(nil))
}

func HashPreImage(api frontend.API, variable frontend.Variable) frontend.Variable {
	m, _ := mimc.NewMiMC(api)
	m.Write(variable)

	return m.Sum()
}
