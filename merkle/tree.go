package merkle_tree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/cbergoon/merkletree"
	gnarkMerkleTree "github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/influxdata/influxdb/pkg/bytesutil"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/sha3"
	"hash"
)

type ByteContent struct {
	B        []byte
	HashFunc func() hash.Hash
}

type Tree struct {
	*merkletree.MerkleTree
	Content []merkletree.Content
	Hashes  [][]byte

	HashFunc func() hash.Hash
	isZKTree bool
}

// CalculateHash hashes the values of a TestContent
func (c ByteContent) CalculateHash() ([]byte, error) {
	h := c.HashFunc()
	if _, err := h.Write(c.B); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (c ByteContent) Equals(other merkletree.Content) (bool, error) {
	return bytes.Equal(c.B, other.(ByteContent).B), nil
}

func newMerkleTree(contents []ByteContent, isZKTree bool) (*Tree, error) {
	hashStrategy := sha3.NewLegacyKeccak256
	if isZKTree {
		hashStrategy = mimc.NewMiMC
	}

	var list []merkletree.Content
	var hashes [][]byte
	for i := range contents {
		list = append(list, contents[i])
		h, err := contents[i].CalculateHash()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, h)
	}

	merkleTree, err := merkletree.NewTreeWithHashStrategySorted(list, hashStrategy, true)
	if err != nil {
		return nil, err
	}

	return &Tree{
		MerkleTree: merkleTree,
		Content:    list,
		Hashes:     hashes,
		HashFunc:   hashStrategy,
		isZKTree:   isZKTree,
	}, nil
}

func NewMerkleTreeBytes(bs [][]byte) (*Tree, error) {
	bytesutil.Sort(bs)

	return newMerkleTree(funk.Map(bs, func(b []byte) ByteContent {
		return ByteContent{B: b, HashFunc: sha3.NewLegacyKeccak256}
	}).([]ByteContent), false)
}

func NewMerkleTreeBytesZK(bs [][]byte) (*Tree, error) {
	bytesutil.Sort(bs)

	return newMerkleTree(funk.Map(bs, func(b []byte) ByteContent {
		return ByteContent{B: b, HashFunc: mimc.NewMiMC}
	}).([]ByteContent), true)
}

func (t *Tree) GetRoot() []byte {
	return t.MerkleRoot()
}

func (t *Tree) NumLeaves() int {
	return len(t.Content)
}

func (t *Tree) RootHex() string {
	return hex.EncodeToString(t.MerkleRoot())
}

func (t *Tree) GetLeaf(c ByteContent) ([]byte, error) {
	if c.HashFunc == nil {
		c.HashFunc = t.HashFunc
	}
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
	return t.GetProof(ByteContent{B: b})
}

func (t *Tree) GetHexProofByByte(b []byte) ([]string, error) {
	return t.GetHexProof(ByteContent{B: b})
}

func (t *Tree) GetHexProof(c ByteContent) ([]string, error) {
	proofs, err := t.GetProof(c)
	if err != nil {
		return nil, err
	}
	var proofsHex []string
	for proof := range proofs {
		proofsHex = append(proofsHex, fmt.Sprintf("0x%x", hex.EncodeToString(proofs[proof])))
	}

	return proofsHex, err
}

func (t *Tree) BuilderProofFromLeafByte(leafHash []byte) (merkleRoot []byte, merkleProof [][]byte, proofIndex uint64, err error) {
	var buf bytes.Buffer
	for i, h := range t.Hashes {
		buf.Write(h)
		if bytes.Equal(h, leafHash) {
			proofIndex = uint64(i)
		}
	}
	//fmt.Println("proofIndex", proofIndex)
	segmentSize := 32
	merkleRoot, merkleProof, _, err = gnarkMerkleTree.BuildReaderProof(&buf, t.HashFunc(), segmentSize, proofIndex)
	if err != nil {
		return nil, nil, proofIndex, err
	}
	return merkleRoot, merkleProof, proofIndex, nil
}

func (t *Tree) BuilderProofHelper(leafHash []byte) (merkleRoot []byte, merkleProof [][]byte, proofHelper []int, err error) {
	merkleRoot, merkleProof, proofIndex, err := t.BuilderProofFromLeafByte(leafHash)
	if err != nil {
		return nil, nil, nil, err
	}

	proofHelper = merkle.GenerateProofHelper(merkleProof, proofIndex, uint64(len(t.Content)))

	return merkleRoot, merkleProof, proofHelper, nil
}
