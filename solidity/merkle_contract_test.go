package solidity

import (
	"encoding/hex"
	"fmt"
	"github.com/thoas/go-funk"
	merkle "gnark-bid/merkle"
	"log"
	"math/big"
	"testing"
	"unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type ExportSolidityTestSuiteMerkleTree struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend
	// contract
	merkleContract *MerkleProof
}

func TestRunExportSolidityTestSuiteMerkleTree(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteMerkleTree))
}

func (t *ExportSolidityTestSuiteMerkleTree) SetupTest() {

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
	_, _, m, err := DeployMerkleProof(auth, t.backend)
	t.NoError(err, "deploy contract failed")

	t.merkleContract = m
	t.backend.Commit()
}

func (t *ExportSolidityTestSuiteMerkleTree) TestVerifyProof() {
	var list [][]byte
	for i := 0; i < 16; i++ {
		s := fmt.Sprintf("%dabc", i+1)
		fmt.Println("s", s)
		list = append(list, []byte(s))
	}

	//Create a new Merkle Tree from the list of ByteContent
	tree, err := merkle.NewMerkleTreeBytes(list)
	t.NoError(err, "create merkle tree failed")
	//fmt.Println("all hex leaves:", tree.AllHexLeaves())
	//Get the Merkle Root of the tree
	log.Println("Merkle Root: ", tree.RootHex())

	//String representation
	hash, _ := tree.Content[0].CalculateHash()
	log.Println("hash leaf:", hex.EncodeToString(hash))

	proof, err := tree.GetProof(tree.Content[0])
	t.NoError(err, "get proof hex failed")

	log.Println(proof)
	verified, err := t.merkleContract.Verify(nil, byte32(tree.GetRoot()), byte32(hash), funk.Map(proof, func(s []byte) [32]byte {
		return byte32(s)
	}).([][32]byte))

	t.NoError(err, "call verify failed")
	if t.True(verified, "verify failed") {
		log.Println("merkle tree verified")
	}
}

func byte32(s []byte) (a [32]byte) {
	if len(a) <= len(s) {
		a = *(*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
}
