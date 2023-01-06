package zk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

const MerkleTreeDepth = 10

const fpSize = 4 * 8

func ParserProof(proofBytes []byte) *Proof {
	proof := &Proof{}
	proof.A[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	proof.A[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	proof.B[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	proof.B[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	proof.B[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	proof.B[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	proof.C[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	proof.C[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])
	//fmt.Println("a", proof.A)
	//fmt.Println("b", proof.B)
	//fmt.Println("c", proof.C)

	return proof
}

func ProofToBytes(proof *Proof) []byte {
	proofBytes := make([]byte, 0, fpSize*8)
	proofBytes = append(proofBytes, proof.A[0].Bytes()...)
	proofBytes = append(proofBytes, proof.A[1].Bytes()...)
	proofBytes = append(proofBytes, proof.B[0][0].Bytes()...)
	proofBytes = append(proofBytes, proof.B[0][1].Bytes()...)
	proofBytes = append(proofBytes, proof.B[1][0].Bytes()...)
	proofBytes = append(proofBytes, proof.B[1][1].Bytes()...)
	proofBytes = append(proofBytes, proof.C[0].Bytes()...)
	proofBytes = append(proofBytes, proof.C[1].Bytes()...)

	return proofBytes
}

func FormatVPKey(vkKey string, pkKey string) (*VPKey, error) {
	vpKey := &VPKey{}

	vpKey.PK = groth16.NewProvingKey(ecc.BN254)
	{
		start := time.Now().Unix()
		fmt.Println("start read pk", start)
		pkBuf := bytes.NewBuffer(common.FromHex(pkKey))
		fmt.Println("NewBuffer:", time.Now().Unix()-start)
		if _, err := vpKey.PK.UnsafeReadFrom(pkBuf); err != nil {
			return nil, err
		}
		fmt.Println("ReadFrom time:", time.Now().Unix()-start)
	}

	vpKey.VK = groth16.NewVerifyingKey(ecc.BN254)
	{
		vkBuf := bytes.NewBuffer(common.FromHex(vkKey))
		if _, err := vpKey.VK.ReadFrom(vkBuf); err != nil {
			return nil, err
		}
	}

	return vpKey, nil
}

func ParseMapVPData(jsonBytes []byte) (map[string]*VPKey, error) {
	var vpKey map[string]*VPKey
	if err := json.Unmarshal(jsonBytes, &vpKey); err != nil {
		return nil, err
	}
	formatVP := make(map[string]*VPKey)
	for k, v := range vpKey {
		vp, err := FormatVPKey(v.VerifyingKey, v.ProvingKey)
		if err != nil {
			return nil, err
		}
		formatVP[k] = &VPKey{
			ProvingKey:   v.ProvingKey,
			VerifyingKey: v.VerifyingKey,
			VK:           vp.VK,
			PK:           vp.PK,
		}
	}

	return formatVP, nil
}
