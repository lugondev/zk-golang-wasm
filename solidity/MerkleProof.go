// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MerkleProofMetaData contains all meta data concerning the MerkleProof contract.
var MerkleProofMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"leaf\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610491806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80633423e54814610030575b600080fd5b61004a600480360381019061004591906102b9565b610060565b6040516100579190610343565b60405180910390f35b60008083905060005b83518110156101085760008482815181106100875761008661035e565b5b602002602001015190508083116100c85782816040516020016100ab9291906103ae565b6040516020818303038152906040528051906020012092506100f4565b80836040516020016100db9291906103ae565b6040516020818303038152906040528051906020012092505b50808061010090610413565b915050610069565b508481149150509392505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b61013d8161012a565b811461014857600080fd5b50565b60008135905061015a81610134565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6101ae82610165565b810181811067ffffffffffffffff821117156101cd576101cc610176565b5b80604052505050565b60006101e0610116565b90506101ec82826101a5565b919050565b600067ffffffffffffffff82111561020c5761020b610176565b5b602082029050602081019050919050565b600080fd5b6000610235610230846101f1565b6101d6565b905080838252602082019050602084028301858111156102585761025761021d565b5b835b81811015610281578061026d888261014b565b84526020840193505060208101905061025a565b5050509392505050565b600082601f8301126102a05761029f610160565b5b81356102b0848260208601610222565b91505092915050565b6000806000606084860312156102d2576102d1610120565b5b60006102e08682870161014b565b93505060206102f18682870161014b565b925050604084013567ffffffffffffffff81111561031257610311610125565b5b61031e8682870161028b565b9150509250925092565b60008115159050919050565b61033d81610328565b82525050565b60006020820190506103586000830184610334565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b6103a86103a38261012a565b61038d565b82525050565b60006103ba8285610397565b6020820191506103ca8284610397565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600061041e82610409565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036104505761044f6103da565b5b60018201905091905056fea264697066735822122038bf3f4d4f359d23efc7e2bffaf074d0286b36e9ae540468753f7451532e3a0b64736f6c63430008110033",
}

// MerkleProofABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleProofMetaData.ABI instead.
var MerkleProofABI = MerkleProofMetaData.ABI

// MerkleProofBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MerkleProofMetaData.Bin instead.
var MerkleProofBin = MerkleProofMetaData.Bin

// DeployMerkleProof deploys a new Ethereum contract, binding an instance of MerkleProof to it.
func DeployMerkleProof(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MerkleProof, error) {
	parsed, err := MerkleProofMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MerkleProofBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MerkleProof{MerkleProofCaller: MerkleProofCaller{contract: contract}, MerkleProofTransactor: MerkleProofTransactor{contract: contract}, MerkleProofFilterer: MerkleProofFilterer{contract: contract}}, nil
}

// MerkleProof is an auto generated Go binding around an Ethereum contract.
type MerkleProof struct {
	MerkleProofCaller     // Read-only binding to the contract
	MerkleProofTransactor // Write-only binding to the contract
	MerkleProofFilterer   // Log filterer for contract events
}

// MerkleProofCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleProofCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleProofTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleProofFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleProofSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleProofSession struct {
	Contract     *MerkleProof      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MerkleProofCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleProofCallerSession struct {
	Contract *MerkleProofCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// MerkleProofTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleProofTransactorSession struct {
	Contract     *MerkleProofTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MerkleProofRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleProofRaw struct {
	Contract *MerkleProof // Generic contract binding to access the raw methods on
}

// MerkleProofCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleProofCallerRaw struct {
	Contract *MerkleProofCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleProofTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleProofTransactorRaw struct {
	Contract *MerkleProofTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleProof creates a new instance of MerkleProof, bound to a specific deployed contract.
func NewMerkleProof(address common.Address, backend bind.ContractBackend) (*MerkleProof, error) {
	contract, err := bindMerkleProof(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleProof{MerkleProofCaller: MerkleProofCaller{contract: contract}, MerkleProofTransactor: MerkleProofTransactor{contract: contract}, MerkleProofFilterer: MerkleProofFilterer{contract: contract}}, nil
}

// NewMerkleProofCaller creates a new read-only instance of MerkleProof, bound to a specific deployed contract.
func NewMerkleProofCaller(address common.Address, caller bind.ContractCaller) (*MerkleProofCaller, error) {
	contract, err := bindMerkleProof(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleProofCaller{contract: contract}, nil
}

// NewMerkleProofTransactor creates a new write-only instance of MerkleProof, bound to a specific deployed contract.
func NewMerkleProofTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleProofTransactor, error) {
	contract, err := bindMerkleProof(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleProofTransactor{contract: contract}, nil
}

// NewMerkleProofFilterer creates a new log filterer instance of MerkleProof, bound to a specific deployed contract.
func NewMerkleProofFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleProofFilterer, error) {
	contract, err := bindMerkleProof(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleProofFilterer{contract: contract}, nil
}

// bindMerkleProof binds a generic wrapper to an already deployed contract.
func bindMerkleProof(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MerkleProofABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleProof *MerkleProofRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleProof.Contract.MerkleProofCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleProof *MerkleProofRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleProof.Contract.MerkleProofTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleProof *MerkleProofRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleProof.Contract.MerkleProofTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleProof *MerkleProofCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleProof.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleProof *MerkleProofTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleProof.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleProof *MerkleProofTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleProof.Contract.contract.Transact(opts, method, params...)
}

// Verify is a free data retrieval call binding the contract method 0x3423e548.
//
// Solidity: function verify(bytes32 root, bytes32 leaf, bytes32[] proof) pure returns(bool)
func (_MerkleProof *MerkleProofCaller) Verify(opts *bind.CallOpts, root [32]byte, leaf [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _MerkleProof.contract.Call(opts, &out, "verify", root, leaf, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x3423e548.
//
// Solidity: function verify(bytes32 root, bytes32 leaf, bytes32[] proof) pure returns(bool)
func (_MerkleProof *MerkleProofSession) Verify(root [32]byte, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _MerkleProof.Contract.Verify(&_MerkleProof.CallOpts, root, leaf, proof)
}

// Verify is a free data retrieval call binding the contract method 0x3423e548.
//
// Solidity: function verify(bytes32 root, bytes32 leaf, bytes32[] proof) pure returns(bool)
func (_MerkleProof *MerkleProofCallerSession) Verify(root [32]byte, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _MerkleProof.Contract.Verify(&_MerkleProof.CallOpts, root, leaf, proof)
}
