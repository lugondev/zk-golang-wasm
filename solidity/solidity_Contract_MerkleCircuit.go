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

// MerkleCircuitMetaData contains all meta data concerning the MerkleCircuit contract.
var MerkleCircuitMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566050600b82828239805160001a6073146043577f4e487b7100000000000000000000000000000000000000000000000000000000600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212209d4df6224380ed7adc821575cb7e4144c3f25c6e0b00be34fbf4453019beaf1464736f6c63430008110033",
}

// MerkleCircuitABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleCircuitMetaData.ABI instead.
var MerkleCircuitABI = MerkleCircuitMetaData.ABI

// MerkleCircuitBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MerkleCircuitMetaData.Bin instead.
var MerkleCircuitBin = MerkleCircuitMetaData.Bin

// DeployMerkleCircuit deploys a new Ethereum contract, binding an instance of MerkleCircuit to it.
func DeployMerkleCircuit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MerkleCircuit, error) {
	parsed, err := MerkleCircuitMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MerkleCircuitBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MerkleCircuit{MerkleCircuitCaller: MerkleCircuitCaller{contract: contract}, MerkleCircuitTransactor: MerkleCircuitTransactor{contract: contract}, MerkleCircuitFilterer: MerkleCircuitFilterer{contract: contract}}, nil
}

// MerkleCircuit is an auto generated Go binding around an Ethereum contract.
type MerkleCircuit struct {
	MerkleCircuitCaller     // Read-only binding to the contract
	MerkleCircuitTransactor // Write-only binding to the contract
	MerkleCircuitFilterer   // Log filterer for contract events
}

// MerkleCircuitCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleCircuitCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleCircuitTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleCircuitTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleCircuitFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleCircuitFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleCircuitSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleCircuitSession struct {
	Contract     *MerkleCircuit    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MerkleCircuitCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleCircuitCallerSession struct {
	Contract *MerkleCircuitCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// MerkleCircuitTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleCircuitTransactorSession struct {
	Contract     *MerkleCircuitTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// MerkleCircuitRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleCircuitRaw struct {
	Contract *MerkleCircuit // Generic contract binding to access the raw methods on
}

// MerkleCircuitCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleCircuitCallerRaw struct {
	Contract *MerkleCircuitCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleCircuitTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleCircuitTransactorRaw struct {
	Contract *MerkleCircuitTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleCircuit creates a new instance of MerkleCircuit, bound to a specific deployed contract.
func NewMerkleCircuit(address common.Address, backend bind.ContractBackend) (*MerkleCircuit, error) {
	contract, err := bindMerkleCircuit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleCircuit{MerkleCircuitCaller: MerkleCircuitCaller{contract: contract}, MerkleCircuitTransactor: MerkleCircuitTransactor{contract: contract}, MerkleCircuitFilterer: MerkleCircuitFilterer{contract: contract}}, nil
}

// NewMerkleCircuitCaller creates a new read-only instance of MerkleCircuit, bound to a specific deployed contract.
func NewMerkleCircuitCaller(address common.Address, caller bind.ContractCaller) (*MerkleCircuitCaller, error) {
	contract, err := bindMerkleCircuit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleCircuitCaller{contract: contract}, nil
}

// NewMerkleCircuitTransactor creates a new write-only instance of MerkleCircuit, bound to a specific deployed contract.
func NewMerkleCircuitTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleCircuitTransactor, error) {
	contract, err := bindMerkleCircuit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleCircuitTransactor{contract: contract}, nil
}

// NewMerkleCircuitFilterer creates a new log filterer instance of MerkleCircuit, bound to a specific deployed contract.
func NewMerkleCircuitFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleCircuitFilterer, error) {
	contract, err := bindMerkleCircuit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleCircuitFilterer{contract: contract}, nil
}

// bindMerkleCircuit binds a generic wrapper to an already deployed contract.
func bindMerkleCircuit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MerkleCircuitABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleCircuit *MerkleCircuitRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleCircuit.Contract.MerkleCircuitCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleCircuit *MerkleCircuitRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleCircuit.Contract.MerkleCircuitTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleCircuit *MerkleCircuitRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleCircuit.Contract.MerkleCircuitTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleCircuit *MerkleCircuitCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleCircuit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleCircuit *MerkleCircuitTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleCircuit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleCircuit *MerkleCircuitTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleCircuit.Contract.contract.Transact(opts, method, params...)
}
