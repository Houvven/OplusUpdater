package test

import (
	"github.com/deatil/go-cryptobin/cryptobin/rsa"
	"testing"
)

func TestRsaPublicKeyGen(t *testing.T) {

	obj := rsa.New().
		GenerateKey(2048)

	PubKeyPem := obj.
		CreatePublicKey().
		ToKeyString()

	println("public key:", PubKeyPem)
}
