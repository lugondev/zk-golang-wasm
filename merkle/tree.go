package merkle_tree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/cbergoon/merkletree"
	gnarkMerkleTree "github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/influxdata/influxdb/pkg/bytesutil"
	"golang.org/x/crypto/sha3"
)

type ByteContent struct {
	b []byte
}

type Tree struct {
	*merkletree.MerkleTree
	Content []merkletree.Content
	Hashes  [][]byte
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

func newMerkleTree(contents []ByteContent) (*Tree, error) {
	var list []merkletree.Content
	var hashes [][]byte
	for i := range contents {
		list = append(list, contents[i])
		hashes = append(hashes, crypto.Keccak256(contents[i].b))
	}

	merkleTree, err := merkletree.NewTreeWithHashStrategySorted(list, sha3.NewLegacyKeccak256, true)
	if err != nil {
		return nil, err
	}
	return &Tree{
		MerkleTree: merkleTree,
		Content:    list,
		Hashes:     hashes,
	}, nil
}

func NewMerkleTreeBytes(bs [][]byte) (*Tree, error) {
	bytesutil.Sort(bs)
	var list []ByteContent
	for i := range bs {
		list = append(list, ByteContent{b: bs[i]})
	}

	return newMerkleTree(list)
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

func (t *Tree) GetHexProofByByte(b []byte) ([]string, error) {
	return t.GetHexProof(ByteContent{b: b})
}

func (t *Tree) GetHexProof(c merkletree.Content) ([]string, error) {
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

	segmentSize := 32
	merkleRoot, merkleProof, _, err = gnarkMerkleTree.BuildReaderProof(&buf, hash.MIMC_BN254.New(), segmentSize, proofIndex)
	if err != nil {
		return nil, nil, proofIndex, err
	}
	return merkleRoot, merkleProof, proofIndex, nil
}

func (t *Tree) BuilderProof(c []byte) (merkleRoot []byte, merkleProof [][]byte, proofHelper []int, err error) {
	merkleRoot, merkleProof, proofHelper, err = t.BuilderProofHelper(c)
	return merkleRoot, merkleProof, proofHelper, err
}

func (t *Tree) BuilderProofHelper(leafHash []byte) (merkleRoot []byte, merkleProof [][]byte, proofHelper []int, err error) {
	merkleRoot, merkleProof, proofIndex, err := t.BuilderProofFromLeafByte(leafHash)
	if err != nil {
		return nil, nil, nil, err
	}

	proofHelper = merkle.GenerateProofHelper(merkleProof, proofIndex, uint64(len(t.Content)))

	return merkleRoot, merkleProof, proofHelper, nil
}
