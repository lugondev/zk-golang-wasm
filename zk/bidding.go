package zk

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/thoas/go-funk"
	merkleTree "gnark-bid/merkle"
	zkCircuit "gnark-bid/zk/circuits"
	"math/big"
)

type Identity struct {
	Nullifier  *big.Int
	Commitment *big.Int
	Trapdoor   *big.Int
}

type Bidding struct {
	isReady     bool
	RoomID      int
	Username    string
	PrivateCode *big.Int

	Identity Identity

	mkTree *merkleTree.Tree
	g16    *GnarkGroth16
}

func createMerkleTree(list [][]byte) (*merkleTree.Tree, error) {
	maxLeaves := math.BigPow(2, int64(MerkleTreeDepth))
	// generate leaves: users use this to generate their own merkle tree
	leaves := make([][]byte, int(maxLeaves.Int64()))
	for i := 0; i < len(leaves); i++ {
		if len(list) <= i || list[i] == nil || len(list[i]) == 0 {
			leaves[i] = []byte{}
		} else {
			leaves[i] = list[i]
		}
	}
	mkTree, err := merkleTree.NewMerkleTreeBytesZK(leaves)
	if err != nil {
		return nil, err
	}
	return mkTree, nil
}

func fakeListTesting() [][]byte {
	var list [][]byte
	// generate leaves: users use this to generate their own merkle tree
	for i := 0; i < 20; i++ {
		r := fmt.Sprintf("user=username_%d|room=%d", i+1, 1111)
		list = append(list, new(big.Int).SetBytes([]byte(r)).Bytes())
	}
	return list
}

func NewBidding() (*Bidding, error) {
	mkTree, err := createMerkleTree(fakeListTesting())
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return nil, err
	}
	nullifier, err := poseidon.HashBytes(buf)
	if err != nil {
		return nil, err
	}

	vpKey, err := GetVPKey("BiddingCircuit")
	if err != nil {
		return nil, err
	}
	var c zkCircuit.BiddingCircuit
	c.UserMerklePath = make([]frontend.Variable, MerkleTreeDepth+1)
	c.UserMerkleHelper = make([]frontend.Variable, MerkleTreeDepth)

	g16, err := NewGnarkGroth16(vpKey, &c)
	return &Bidding{
		mkTree: mkTree,
		g16:    g16,
		Identity: Identity{
			Nullifier: nullifier,
		},
		isReady: false,
	}, nil
}

func (b *Bidding) InitSession(roomId int, username string, privateCode *big.Int) error {
	b.RoomID = roomId
	b.Username = username
	b.PrivateCode = privateCode

	err := b.generateIdentity()
	if err != nil {
		return err
	}
	b.isReady = true
	return nil
}

func (b *Bidding) getUserID() *big.Int {
	userInfo := fmt.Sprintf("user=%s|room=%d", b.Username, b.RoomID)
	return new(big.Int).SetBytes([]byte(userInfo))
}

func (b *Bidding) getUserLeaf() []byte {
	return zkCircuit.HashMIMC(b.getUserID().Bytes()).Bytes()
}

func (b *Bidding) generateIdentity() error {
	commitment, err := poseidon.Hash([]*big.Int{b.getUserID(), b.PrivateCode})
	if err != nil {
		return err
	}
	trapdoorNumber, err := poseidon.Hash([]*big.Int{commitment, b.Identity.Nullifier})
	if err != nil {
		return err
	}
	b.Identity = Identity{
		Nullifier:  b.Identity.Nullifier,
		Commitment: commitment,
		Trapdoor:   trapdoorNumber,
	}
	return nil
}

func (b *Bidding) RenewSession() error {
	if !b.isReady {
		return fmt.Errorf("session is not ready")
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	var err error
	b.Identity.Nullifier, err = poseidon.HashBytes(buf)
	if err != nil {
		return err
	}

	return b.generateIdentity()
}

func (b *Bidding) GetIdentity() Identity {
	return b.Identity
}

func (b *Bidding) JoinRoom(roomId int) error {
	b.RoomID = roomId
	return b.RenewSession()
}

func (b *Bidding) GetProof(bidValue *big.Int) (*Proof, [5]*big.Int, error) {
	if !b.isReady {
		return nil, [5]*big.Int{}, fmt.Errorf("session is not ready")
	}
	merkleRoot, merkleProof, proofHelper, err := b.mkTree.BuilderProofHelper(b.getUserLeaf())
	if err != nil {
		return nil, [5]*big.Int{}, err
	}

	merkleAssignment := zkCircuit.BiddingCircuit{
		UserMerklePath: funk.Map(merkleProof, func(p []byte) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleHelper: funk.Map(proofHelper, func(p int) frontend.Variable {
			return p
		}).([]frontend.Variable),
		UserMerkleRoot: merkleRoot,
		BidValue:       bidValue,
		Identity: zkCircuit.Identity{
			Nullifier:  b.Identity.Nullifier,
			Commitment: b.Identity.Commitment,
			Trapdoor:   b.Identity.Trapdoor,
		},
		UserData: zkCircuit.UserData{
			UserID:      b.getUserID(),
			PrivateCode: b.PrivateCode,
		},
	}

	proofParser, _, err := b.g16.GenerateProof(&merkleAssignment)
	if err != nil {
		return nil, [5]*big.Int{}, err
	}
	var publicInput [5]*big.Int
	publicInput[0] = new(big.Int).SetBytes(merkleRoot)
	publicInput[1] = b.Identity.Nullifier
	publicInput[2] = b.Identity.Commitment
	publicInput[3] = b.Identity.Trapdoor
	publicInput[4] = bidValue

	return proofParser, publicInput, nil
}

func (b *Bidding) VerifyProof(proof *Proof, inputs [5]*big.Int) (bool, error) {
	if !b.isReady {
		return false, fmt.Errorf("session is not ready")
	}

	merklePath := make([]frontend.Variable, MerkleTreeDepth+1)
	for i := 0; i < MerkleTreeDepth+1; i++ {
		merklePath[i] = big.NewInt(0)
	}
	merkleHelper := make([]frontend.Variable, MerkleTreeDepth)
	for i := 0; i < MerkleTreeDepth; i++ {
		merkleHelper[i] = big.NewInt(0)
	}

	assignment := &zkCircuit.BiddingCircuit{
		UserMerklePath:   merklePath,
		UserMerkleHelper: merkleHelper,
		UserMerkleRoot:   inputs[0],
		BidValue:         inputs[4],
		Identity: zkCircuit.Identity{
			Nullifier:  inputs[1],
			Commitment: inputs[2],
			Trapdoor:   inputs[3],
		},
		UserData: zkCircuit.UserData{
			UserID:      big.NewInt(0),
			PrivateCode: big.NewInt(0),
		},
	}

	g16Proof := groth16.NewProof(ecc.BN254)
	proofBuf := bytes.NewBuffer(ProofToBytes(proof))
	fmt.Println("compare:", bytes.Equal(ProofToBytes(proof), proofBuf.Bytes()))

	if _, err := g16Proof.ReadFrom(proofBuf); err != nil {
		fmt.Println("read proof error", err)
		return false, err
	}

	return b.g16.VerifyProof(assignment, g16Proof)
}
