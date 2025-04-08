package updater

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/tidwall/pretty"
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

func (r *ResponseResult) PrettyPrint() {
	var body map[string]interface{}
	_ = json.Unmarshal(r.DecryptedBodyBytes, &body)

	m := map[string]interface{}{
		"responseCode": r.ResponseCode,
		"errMsg":       r.ErrMsg,
		"body":         body,
	}

	if bytes, err := json.Marshal(m); err == nil {
		fmt.Println(string(pretty.Color(pretty.Pretty(bytes), nil)))
	}
}
