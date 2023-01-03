package main

import (
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/thoas/go-funk"
	"gnark-bid/zk"
	"gnark-bid/zk/circuits"
	"log"
	"reflect"
	"strings"
)

type VP struct {
	Name string
	Key  *zk.VPKey
}

func main() {
	var cBid zk_circuit.BidCircuit
	cBid.UserMerklePath = make([]frontend.Variable, zk.MerkleTreeDepth+1)
	cBid.UserMerkleHelper = make([]frontend.Variable, zk.MerkleTreeDepth)
	listCircuit := []frontend.Circuit{
		&cBid,
	}

	keys := funk.Map(listCircuit, func(circuit frontend.Circuit) VP {
		name := reflect.TypeOf(circuit).String()
		structName := lastString(strings.Split(name, "."))
		fmt.Println("circuit initializing:", structName)
		k, err := zk.InitGroth16(circuit, structName)
		if err != nil {
			log.Fatal("groth16 error:", err)
		}
		return VP{
			Name: structName,
			Key:  k,
		}
	}).([]VP)

	mapKeys := make(map[string]*zk.VPKey)
	for _, k := range keys {
		mapKeys[k.Name] = k.Key
	}

	if err := zk.WriteTemplateKey(mapKeys); err != nil {
		log.Fatal("write template error:", err)
	}
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}
