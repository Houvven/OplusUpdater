package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
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

func TestDeviceIdGen(t *testing.T) {
	id := updater.GenerateDeviceId("864290073152698")
	fmt.Printf("device id: %s, len: %d\n", id, len(id))
}
