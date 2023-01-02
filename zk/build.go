package zk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"text/template"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func InitGroth16(c frontend.Circuit, name string) (*VPKey, error) {
	r1csCompiled, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, c)
	if err != nil {
		return nil, err
	}

	pk, vk, err := groth16.Setup(r1csCompiled)
	if err != nil {
		return nil, err
	}

	//{
	//	f, err := os.Create("solidity/zk.g16.vk")
	//	if err != nil {
	//		return err
	//	}
	//	_, err = vk.WriteRawTo(f)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//{
	//	f, err := os.Create("solidity/zk.g16.pk")
	//	if err != nil {
	//		return err
	//	}
	//	_, err = pk.WriteRawTo(f)
	//	if err != nil {
	//		return err
	//	}
	//}

	{
		f, err := os.Create(fmt.Sprintf("solidity/Contract_%s.sol", name))
		if err != nil {
			return nil, err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return nil, err
		}
	}
	return GenerateKey(pk, vk)
}

func GenerateKey(pk groth16.ProvingKey, vk groth16.VerifyingKey) (*VPKey, error) {
	bufVk := new(bytes.Buffer)
	bufPk := new(bytes.Buffer)
	_, err := vk.WriteRawTo(bufVk)
	if err != nil {
		return nil, err
	}
	_, err = pk.WriteRawTo(bufPk)
	if err != nil {
		return nil, err
	}

	return &VPKey{
		ProvingKey:   common.Bytes2Hex(bufPk.Bytes()),
		VerifyingKey: common.Bytes2Hex(bufVk.Bytes()),
	}, nil
}

func WriteJsonFile(jsonBytes []byte, fileName string) error {
	bufJson := new(bytes.Buffer)
	bufJson.Write(jsonBytes)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = bufJson.WriteTo(f)
	return err
}

func WriteTemplateKey(data map[string]*VPKey) error {
	jsonBytes, _ := json.Marshal(data)
	return GenerateTemplateKey(jsonBytes, "wasm/parse_vp_key.go")
}

func GenerateTemplateKey(jsonBytes []byte, fileName string) error {
	getwd, err := os.Getwd()
	if err != nil {
		return err
	}
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/zk/build/parse_vp_key-template.go.temp", getwd)))
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, string(jsonBytes)); err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = tpl.WriteTo(f)
	return err
}
