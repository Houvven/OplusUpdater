package util

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"github.com/deatil/go-cryptobin/cryptobin/rsa"
	"io"
	"strconv"
	"strings"
	"time"
)

func RandomIv() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	return iv, err
}

func RandomKey() ([]byte, error) {
	key := make([]byte, aes.BlockSize*2)
	_, err := io.ReadFull(rand.Reader, key)
	return key, err
}

func GenerateProtectedVersion() string {
	return strconv.FormatInt(time.Now().Add(time.Hour*24).UnixNano(), 10)
}

func GenerateProtectedKey(key []byte, pubKey []byte) (string, error) {
	encrypt := rsa.New().
		FromString(base64.StdEncoding.EncodeToString(key)).
		FromPublicKey(pubKey).
		EncryptOAEP()
	return encrypt.ToBase64String(), encrypt.Error()
}

func GetDefaultDeviceId() string {
	return strings.Repeat("0", 64)
}
