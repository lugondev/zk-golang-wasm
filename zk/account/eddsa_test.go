package account

import (
	"crypto/rand"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func Test_EDDSA(t *testing.T) {
	hFunc := hash.MIMC_BN254.New()

	// create a eddsa key pair
	privateKey, _ := eddsa.New(twistededwards.BN254, rand.Reader)
	publicKey := privateKey.Public()
	t.Log("public key", common.Bytes2Hex(publicKey.Bytes()))

	// note that the message is on 4 bytes
	msg := []byte{0xde, 0xad, 0xf0, 0x0d}

	// sign the message
	signature, _ := privateKey.Sign(msg, hFunc)

	// verifies signature
	isValid, _ := publicKey.Verify(signature, msg, hFunc)

	if !isValid {
		fmt.Println("1. invalid signature")
	} else {
		fmt.Println("1. valid signature")
	}
}
