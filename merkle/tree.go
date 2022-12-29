package merkle

import (
	"bytes"
	"fmt"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type ByteContent struct {
	b []byte
}

type Tree struct {
	*merkletree.MerkleTree
	content []merkletree.Content
}

// CalculateHash hashes the values of a TestContent
func (c ByteContent) CalculateHash() ([]byte, error) {
	h := sha3.NewLegacyKeccak256()
	if _, err := h.Write(c.b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (c ByteContent) Equals(other merkletree.Content) (bool, error) {
	return bytes.Equal(c.b, other.(ByteContent).b), nil
}

func NewMerkleTree(contents []ByteContent) (*Tree, error) {
	var list []merkletree.Content
	for i := range contents {
		list = append(list, contents[i])
	}

	merkleTree, err := merkletree.NewTreeWithHashStrategy(list, sha3.NewLegacyKeccak256)
	if err != nil {
		return nil, err
	}
	return &Tree{
		MerkleTree: merkleTree,
		content:    list,
	}, nil
}

func NewMerkleTreeStrings(strings []string) (*Tree, error) {
	var list []ByteContent
	for i := range strings {
		list = append(list, ByteContent{b: []byte(strings[i])})
	}

	return NewMerkleTree(list)
}

func NewMerkleTreeBytes(bs [][]byte) (*Tree, error) {
	var list []ByteContent
	for i := range bs {
		list = append(list, ByteContent{b: bs[i]})
	}

	return NewMerkleTree(list)
}

func (t *Tree) GetLeaf(c merkletree.Content) ([]byte, error) {
	return c.CalculateHash()
}

func (t *Tree) GetProof(c merkletree.Content) ([][]byte, error) {
	merklePath, _, err := t.GetMerklePath(c)
	if err != nil {
		return nil, err
	}
	var proofs [][]byte
	for path := range merklePath {
		proofs = append(proofs, merklePath[path])
	}
	return proofs, err
}

func (t *Tree) GetProofByByte(b []byte) ([][]byte, error) {
	return t.GetProof(ByteContent{b: b})
}

func (t *Tree) GetProofByByteHex(b []byte) ([]string, error) {
	return t.GetProofHex(ByteContent{b: b})
}

func (t *Tree) GetProofHex(c merkletree.Content) ([]string, error) {
	proofs, err := t.GetProof(c)
	if err != nil {
		return nil, err
	}
	var proofsHex []string
	for proof := range proofs {
		proofsHex = append(proofsHex, fmt.Sprintf("0x%x", common.Bytes2Hex(proofs[proof])))
	}

	return proofsHex, err
}
