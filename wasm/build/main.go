package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"gnark-bid/wasm"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"math/big"
	"syscall/js"
)

var validation = validator.New()

func jsErr(err error, message string) string {
	if message == "" {
		return fmt.Sprintf("{'error': '%s','message': '%s'}", err.Error(), message)
	}
	if err == nil {
		return fmt.Sprintf("{'error': '%s'}", message)
	}
	return fmt.Sprintf("{'error': '%s'}", err.Error())
}

func hash() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		fmt.Println("args", args)

		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}

		value := common.FromHex(args[0].String())
		hashMIMC := zk_circuit.HashMIMC(value)
		return fmt.Sprintf("{'numHash': '%s','hexHash': '%s'}", hashMIMC.String(), hashMIMC.Text(16))
	})
}

func generateProof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}

		inputBytes := common.FromHex(args[0].String())
		privateValue := new(big.Int).SetBytes(inputBytes)
		fmt.Println("privateValue", privateValue.String())

		assignment := zk_circuit.PrivateValueCircuit{
			PrivateValue: privateValue.String(),
			Hash:         zk_circuit.HashMIMC(inputBytes).String(),
		}

		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}

		vkKey, err := wasm.GetVPKey("PrivateValueCircuit")
		if err != nil {
			return jsErr(err, "Cannot read keys")
		}

		var c zk_circuit.PrivateValueCircuit
		g16, err := zk.NewGnarkGroth16(vkKey, &c)
		if err != nil {
			return jsErr(err, "")
		}
		inputProof := [1]*big.Int{zk_circuit.HashMIMC(inputBytes)}
		proofGenerated, err := g16.GenerateProof(assignment, inputProof)
		if err != nil {
			return jsErr(err, "Cannot generate proof")
		}
		proofJSON, _ := json.Marshal(proofGenerated)

		return string(proofJSON)
	})
}

func verifyProof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		var proof zk.Proof
		if err := json.Unmarshal([]byte(args[0].String()), &proof); err != nil {
			return jsErr(err, "Invalid proof passed")
		}
		if err := validation.Var(args[1].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}

		hashBytes := common.FromHex(args[1].String())
		hashValue := new(big.Int).SetBytes(hashBytes)
		fmt.Println("hash", hashValue.String())

		assignment := zk_circuit.PrivateValueCircuit{
			PrivateValue: hashValue,
			Hash:         hashValue.String(),
		}

		vkKey, err := wasm.GetVPKey("PrivateValueCircuit")
		if err != nil {
			return jsErr(err, "Cannot read keys")
		}

		witness, err := frontend.NewWitness(&assignment, ecc.BN254)
		if err != nil {
			return jsErr(err, "Cannot create witness")
		}
		publicWitness, _ := witness.Public()
		proofBytes := zk.ProofToBytes(proof)
		proofG16 := groth16.NewProof(ecc.BN254)
		if _, err := proofG16.ReadFrom(bytes.NewReader(proofBytes)); err != nil {
			return jsErr(err, "Cannot read proof")
		}

		if err := groth16.Verify(proofG16, vkKey.VK, publicWitness); err != nil {
			return jsErr(err, "Cannot verify proof")
		} else {
			return jsErr(nil, "Proof verified")
		}
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("hash", hash())
	js.Global().Set("generateProof", generateProof())
	js.Global().Set("verifyProof", verifyProof())
	<-make(chan bool)
}
