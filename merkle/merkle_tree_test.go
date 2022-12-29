package merkle_test

import (
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
	"log"
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

func Test_MerkleTree(t *testing.T) {
	//Build list of ByteContent to build tree
	var list []merkletree.Content
	list = append(list, TestContent{x: "a"})
	list = append(list, TestContent{x: "b"})
	list = append(list, TestContent{x: "c"})
	list = append(list, TestContent{x: "d"})

	//Create a new Merkle Tree from the list of ByteContent
	tree, err := merkletree.NewTreeWithHashStrategy(list, sha3.NewLegacyKeccak256)
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

	//Verify a specific content in in the tree
	vc, err := tree.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Verify ByteContent: ", vc)

	//String representation
	log.Println(tree)
	hash, err := list[0].CalculateHash()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("hash leaf:", common.Bytes2Hex(hash))
	merklePath, i, err := tree.GetMerklePath(list[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i)
	log.Println(merklePath)
	for path := range merklePath {
		log.Println(common.Bytes2Hex(merklePath[path]))
	}
	//log.Println(string(marshal))
}
