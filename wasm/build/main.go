package main

import (
	"fmt"
	"gnark-bid/zk"
	"syscall/js"
)

var bidding *zk.Bidding

func initBidding() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {

		b, err := zk.NewBidding()
		if err != nil {
			return jsErr(err, "Cannot init bidding")
		}
		bidding = b
		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Session initialized")
	})
}

func initSession() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding == nil {
			return jsErr(nil, "Session not initialized")
		}

		if len(args) != 3 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		username := args[0].String()
		roomID := args[1].Int()
		privateCode, errStr := parserHexToBigInt(args[2].String())
		if errStr != "" {
			return errStr
		}

		if err := bidding.InitSession(roomID, username, privateCode); err != nil {
			return jsErr(err, "Cannot init session")
		}

		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Session initialized")
	})
}

func renewSession(bidding *zk.Bidding) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		err := bidding.RenewSession()
		if err != nil {
			return jsErr(err, "Cannot renew session")
		}
		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Session renewed")
	})
}

func generateProof(bidding *zk.Bidding) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}
		return ""
		//inputBytes := common.FromHex(args[0].String())
		//privateValue := new(big.Int).SetBytes(inputBytes)
		//fmt.Println("privateValue", privateValue.String())
		//
		//assignment := zk_circuit.PrivateValueCircuit{
		//	PrivateValue: privateValue.String(),
		//	Hash:         zk_circuit.HashMIMC(inputBytes).String(),
		//}
		//
		//if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
		//	return jsErr(err, "Invalid argument input passed")
		//}
		//
		//vkKey, err := zk.GetVPKey("PrivateValueCircuit")
		//if err != nil {
		//	return jsErr(err, "Cannot read keys")
		//}
		//
		//var c zk_circuit.PrivateValueCircuit
		//g16, err := zk.NewGnarkGroth16(vkKey, &c)
		//if err != nil {
		//	return jsErr(err, "")
		//}
		//
		//inputProof := []*big.Int{zk_circuit.HashMIMC(inputBytes)}
		//proofGenerated, _, err := g16.GenerateProof(&assignment)
		//if err != nil {
		//	return jsErr(err, "Cannot generate proof")
		//}
		//data := map[string]interface{}{
		//	"proof": proofGenerated,
		//	"input": inputProof,
		//}
		//dataJSON, _ := json.Marshal(data)
		//
		//return string(dataJSON)
	})
}

func verifyProof(bidding *zk.Bidding) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return jsErr(nil, "Invalid no of arguments passed")
		}

		return ""
		//var proof zk.Proof
		//if err := json.Unmarshal([]byte(args[0].String()), &proof); err != nil {
		//	return jsErr(err, "Invalid proof passed")
		//}
		//if err := validation.Var(args[1].String(), "required,hexadecimal"); err != nil {
		//	return jsErr(err, "Invalid argument input passed")
		//}
		//
		//hashBytes := common.FromHex(args[1].String())
		//hashValue := new(big.Int).SetBytes(hashBytes)
		//fmt.Println("hash", hashValue.String())
		//
		//assignment := zk_circuit.PrivateValueCircuit{
		//	PrivateValue: hashValue,
		//	Hash:         hashValue.String(),
		//}
		//
		//vkKey, err := zk.GetVPKey("PrivateValueCircuit")
		//if err != nil {
		//	return jsErr(err, "Cannot read keys")
		//}
		//
		//witness, err := frontend.NewWitness(&assignment, ecc.BN254)
		//if err != nil {
		//	return jsErr(err, "Cannot create witness")
		//}
		//publicWitness, _ := witness.Public()
		//proofBytes := zk.ProofToBytes(proof)
		//proofG16 := groth16.NewProof(ecc.BN254)
		//if _, err := proofG16.ReadFrom(bytes.NewReader(proofBytes)); err != nil {
		//	return jsErr(err, "Cannot read proof")
		//}
		//
		//if err := groth16.Verify(proofG16, vkKey.VK, publicWitness); err != nil {
		//	return jsErr(err, "Cannot verify proof")
		//} else {
		//	return jsErr(nil, "Proof verified")
		//}
	})
}

func main() {

	fmt.Println("Go Web Assembly - Bidding Platform")
	js.Global().Set("initBidding", initBidding())
	js.Global().Set("initSession", initSession())
	js.Global().Set("renewSession", renewSession(bidding))
	js.Global().Set("generateProof", generateProof(bidding))
	js.Global().Set("verifyProof", verifyProof(bidding))
	<-make(chan bool)
}
