package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gnark-bid/zk"
	"math/big"
	"syscall/js"
)

var bidding *zk.Bidding

func initBidding() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding != nil {
			return jsErr(nil, "Session already initialized")
		}

		fmt.Println("Initializing bidding")
		if b, err := zk.NewBidding(nil); err != nil {
			return jsErr(err, "Cannot init bidding")
		} else {
			bidding = b
		}

		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Session initialized")
	})
}

func createSession() js.Func {
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

func renewSession() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding == nil {
			return jsErr(nil, "Session not initialized")
		}

		err := bidding.RenewSession()
		if err != nil {
			return jsErr(err, "Cannot renew session")
		}
		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Session renewed")
	})
}

func isInitialized() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return bidding != nil
	})
}

func generateProof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding == nil {
			return jsErr(nil, "Session not initialized")
		}

		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}
		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			return jsErr(err, "Invalid argument input passed")
		}

		inputBytes := common.FromHex(args[0].String())
		bidValue := new(big.Int).SetBytes(inputBytes)
		proofs, inputs, err := bidding.GetProof(bidValue)
		if err != nil {
			return jsErr(err, "Cannot generate proof")
		}

		data := map[string]interface{}{
			"proofs": proofs,
			"inputs": inputs,
		}
		dataJSON, _ := json.Marshal(data)

		return string(dataJSON)
	})
}

func joinRoom() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding == nil {
			return jsErr(nil, "Session not initialized")
		}

		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}

		if err := bidding.JoinRoom(args[0].Int()); err != nil {
			return jsErr(err, "Cannot join room")
		}

		return fmt.Sprintf("{'status': '%s','message': '%s'}", "success", "Room joined")
	})
}

func getCurrentSession() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if bidding == nil {
			return jsErr(nil, "Session not initialized")
		}

		data := map[string]interface{}{
			"identity":    bidding.GetIdentity(),
			"roomID":      bidding.RoomID,
			"privateCode": bidding.PrivateCode,
			"username":    bidding.Username,
		}
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return jsErr(err, "Cannot marshal data")
		}

		return string(dataJSON)
	})
}

func verifyProof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return jsErr(nil, "Invalid no of arguments passed")
		}

		var data struct {
			Proofs *zk.Proof   `json:"proofs"`
			Inputs [5]*big.Int `json:"inputs"`
		}

		if err := json.Unmarshal([]byte(args[0].String()), &data); err != nil {
			return jsErr(err, "Cannot unmarshal data")
		}

		verified, err := bidding.VerifyProof(data.Proofs, data.Inputs)
		if err != nil {
			return jsErr(err, "Cannot verify proof")
		}
		return verified
	})
}

func main() {
	fmt.Println("Go Web Assembly - Bidding Platform")

	js.Global().Set("isInitialized", isInitialized())
	js.Global().Set("initBidding", initBidding())
	js.Global().Set("createSession", createSession())
	js.Global().Set("renewSession", renewSession())
	js.Global().Set("getCurrentSession", getCurrentSession())
	js.Global().Set("joinRoom", joinRoom())
	js.Global().Set("generateProof", generateProof())
	js.Global().Set("verifyProof", verifyProof())

	<-make(chan bool)
}
