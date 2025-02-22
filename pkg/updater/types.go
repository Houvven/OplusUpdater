package updater

import (
	"encoding/base64"
	"encoding/json"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

type Region = string

const (
	RegionCn = "CN"
	RegionEu = "EU"
	RegionIn = "IN"
	RegionSg = "SG"
)

type CryptoConfig struct {
	ProtectedKey       string `json:"protectedKey"`
	Version            string `json:"version"`
	NegotiationVersion string `json:"negotiationVersion"`
}

type RequestBody struct {
	Cipher string `json:"cipher"`
	Iv     string `json:"iv"`
}

type ResponseResult struct {
	ResponseCode       int    `json:"responseCode"`
	ErrMsg             string `json:"errMsg"`
	Body               any    `json:"body"`
	DecryptedBodyBytes []byte
}

func (r *ResponseResult) DecryptBody(key []byte) error {
	var m map[string]interface{}
	if r.Body == nil {
		return nil
	}
	if err := json.Unmarshal([]byte(r.Body.(string)), &m); err != nil {
		return err
	}

	iv, err := base64.StdEncoding.DecodeString(m["iv"].(string))
	if err != nil {
		return err
	}
	cipherBytes := crypto.FromBase64String(m["cipher"].(string)).
		Aes().CTR().NoPadding().
		WithKey(key).WithIv(iv).
		Decrypt().
		ToBytes()

	r.DecryptedBodyBytes = cipherBytes
	return nil
}
