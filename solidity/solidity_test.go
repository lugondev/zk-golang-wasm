package solidity

import (
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	"gnark-bid/zk"
	"math/big"
)

type ExportSolidityTestSuite[T any] struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend

	// contract
	contract *T

	// groth16 gnark objects
	vk  groth16.VerifyingKey
	pk  groth16.ProvingKey
	g16 *zk.GnarkGroth16
}

func InitSetup[T any](t *ExportSolidityTestSuite[T], circuit frontend.Circuit, deployFunc func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *T, error), vpName string) {
	const gasLimit uint64 = 4712388
	// setup simulated backend
	key, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	t.NoError(err, "init keyed transactor")

	genesis := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(1000000000000000000)}, // 1 Eth
	}
	t.backend = backends.NewSimulatedBackend(genesis, gasLimit)

	// deploy verifier contract
	_, _, v, err := deployFunc(auth, t.backend)
	t.NoError(err, "deploy verifier contract failed")
	t.contract = v
	t.backend.Commit()

	vpKey, err := zk.GetVPKey(vpName)
	t.NoError(err, "getting vpkey failed")

	// read proving and verifying keys
	t.pk = vpKey.PK
	t.vk = vpKey.VK

	// create gnark groth16
	t.g16, err = zk.NewGnarkGroth16(vpKey, circuit)
	t.NoError(err, "init groth16 failed")
}
