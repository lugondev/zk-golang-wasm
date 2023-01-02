package merkle_tree_test

import (
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"testing"
)

type TestContent struct {
	x string
}

// CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha3.NewLegacyKeccak256()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

func TestNewMerkleTree(t *testing.T) {
	//Build list of Content to build tree
	var list []merkletree.Content
	maxLeaves := math.BigPow(2, 10)
	//maxLeaves := big.NewInt(20)
	fmt.Println("maxLeaves:", maxLeaves.String())
	for i := big.NewInt(0); i.Cmp(maxLeaves) < 0; i.Add(i, big.NewInt(1)) {
		//id := i.Int64() + 1
		//fmt.Println("id:", id)
		r := fmt.Sprintf("%dabc", 1)
		list = append(list, TestContent{x: r})
	}
	//for i := 0; i < 100; i++ {
	//	r := fmt.Sprintf("%dabc", i+1)
	//	list = append(list, TestContent{x: r})
	//}
	//Create a new Merkle Tree from the list of Content
	tree, err := merkletree.NewTreeWithHashStrategySorted(list, sha3.NewLegacyKeccak256, true)
	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := tree.MerkleRoot()
	log.Println("Merkle Root: ", common.Bytes2Hex(mr))

	//Verify the entire tree (hashes for each node) is valid
	vt, err := tree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Tree: ", vt)

	//Verify a specific content in the tree
	vc, err := tree.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Verify Content: ", vc)

	//String representation
	c := list[1]
	cHash, _ := c.CalculateHash()
	log.Println("cHash: ", common.Bytes2Hex(cHash))
	merklePath, _, err := tree.GetMerklePath(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Merkle Path len: ", len(merklePath))
	//for _, p := range merklePath {
	//	log.Println(common.Bytes2Hex(p))
	//}
}
