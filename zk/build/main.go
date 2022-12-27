package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gnark-bid/zk"
	"log"
	"os"
	"text/template"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	err := generateGroth16()
	if err != nil {
		log.Fatal("groth16 error:", err)
	}
}

func generateGroth16() error {
	var circuit zk.BidCircuit

	r1csCompiled, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		return err
	}

	pk, vk, err := groth16.Setup(r1csCompiled)
	if err != nil {
		return err
	}
	if err := generateKeys(pk, vk); err != nil {
		return err
	}

	{
		f, err := os.Create("solidity/zk.g16.vk")
		if err != nil {
			return err
		}
		_, err = vk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}

	{
		f, err := os.Create("solidity/zk.g16.pk")
		if err != nil {
			return err
		}
		_, err = pk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}

	{
		f, err := os.Create("solidity/contract_g16.sol")
		if err != nil {
			return err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateKeys(pk groth16.ProvingKey, vk groth16.VerifyingKey) error {
	bufVk := new(bytes.Buffer)
	bufPk := new(bytes.Buffer)
	_, err := vk.WriteRawTo(bufVk)
	if err != nil {
		return err
	}
	_, err = pk.WriteRawTo(bufPk)
	if err != nil {
		return err
	}

	jsonStruct := zk.VPKey{
		ProvingKey:   common.Bytes2Hex(bufPk.Bytes()),
		VerifyingKey: common.Bytes2Hex(bufVk.Bytes()),
	}

	bufJson := new(bytes.Buffer)
	jsonBytes, _ := json.Marshal(jsonStruct)
	if err := generateTemplateKey(jsonBytes); err != nil {
		return err
	}

	bufJson.Write(jsonBytes)
	f, err := os.Create("solidity/zk.json")
	if err != nil {
		return err
	}
	_, err = bufJson.WriteTo(f)

	return err
}

func generateTemplateKey(jsonBytes []byte) error {
	getwd, err := os.Getwd()
	if err != nil {
		return err
	}
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/zk/build/parse_vp_key-template.go", getwd)))
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, string(jsonBytes)); err != nil {
		return err
	}
	f, err := os.Create("wasm/parse_vp_key.go")
	if err != nil {
		return err
	}
	_, err = tpl.WriteTo(f)
	return err
}
