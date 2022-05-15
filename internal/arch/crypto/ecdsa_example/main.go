// Based on https://pkg.go.dev/crypto/ecdsa#example-package

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func signAndVerify(privateKey *ecdsa.PrivateKey) {
	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("signature verified:", valid)
}

func main() {
	{
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\n********* generatedPrivateKey: \n%+v\n", privateKey)
		signAndVerify(privateKey)
	}

	{
		x, _ := new(big.Int).SetString("5479062030798703207418524118268733530497083089004144624114120854136991155223", 0)
		y, _ := new(big.Int).SetString("2001412325671319434120395853596552755969429529316841246921487170062309606243", 0)
		d, _ := new(big.Int).SetString("8905450382323324371539459134339773908523793093040739418479598229768365864881", 0)
		privateKey := ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     x,
				Y:     y,
			},
			D: d,
		}
		fmt.Printf("\n********* givenPrivateKey: \n%+v\n", privateKey)
		signAndVerify(&privateKey)
	}
}
