package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (header *UpdateRequestHeaders) makeHeader() (http.Header, error) {
	h := http.Header{
		"androidVersion": {header.AndroidVersion},
		"colorOSVersion": {header.ColorOSVersion},
		"otaVersion":     {header.OtaVersion},
		"deviceId":       {header.DeviceId},
		"model":          {strings.Split(header.OtaVersion, "_")[0]},
		"language":       {"zh-CN"},
		"nvCarrier":      {"10010111"},
		"version":        {"2"},
		"Content-Type":   {"application/json; charset=utf-8"},
	}

	keyJson, err := json.Marshal(header.ProtectedKey)
	if err != nil {
		return nil, err
	}
	h.Set("protectedKey", string(keyJson))
	return h, nil
}

func makeBody(key, iv []byte, cipher RequestCipher) ([]byte, error) {
	marshal, err := json.Marshal(cipher)
	paramsJson, err := json.Marshal(
		map[string]string{
			"cipher": crypto.FromBytes(marshal).
				Aes().CTR().NoPadding().
				WithKey(key).
				WithIv(iv).
				Encrypt().
				ToBase64String(),
			"iv": base64.StdEncoding.EncodeToString(iv),
		},
	)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]string{"params": string(paramsJson)})
}

func RequestUpdate(
	key, iv []byte,
	updateHeaders UpdateRequestHeaders,
	cipher RequestCipher,
) (ResponseResult, error) {
	body, err := makeBody(key, iv, cipher)
	if err != nil {
		return ResponseResult{}, err
	}
	header, err := updateHeaders.makeHeader()
	if err != nil {
		return ResponseResult{}, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: Host, Path: "/update/v3"},
		Header: header,
		Body:   io.NopCloser(bytes.NewBuffer(body)),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponseResult{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}(resp.Body)

	var result ResponseResult
	if json.NewDecoder(resp.Body).Decode(&result) != nil {
		return ResponseResult{}, err
	}
	return result, nil
}

func UpdateResponseParse(body ResponseResult, key []byte) {
	if body.ResponseCode != 200 {
		log.Fatalf("Update failed: %s", body.ErrMsg)
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(body.Body.(string)), &m); err != nil {
		log.Fatalf("Error unmarshalling response body: %v", err)
	}

	iv, err := base64.StdEncoding.DecodeString(m["iv"].(string))
	if err != nil {
		log.Fatalf("Error decoding IV: %v", err)
	}
	cipherBytes := crypto.FromBase64String(m["cipher"].(string)).
		Aes().CTR().NoPadding().
		WithKey(key).
		WithIv(iv).
		Decrypt().ToBytes()

	var cipher UpdateResponseCipher
	if err := json.Unmarshal(cipherBytes, &cipher); err != nil {
		log.Fatalf("Error unmarshalling cipher: %v", err)
	}

	cipherJson, err := json.MarshalIndent(cipher, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting cipher: %v", err)
	}
	log.Printf("Update response: %s", cipherJson)
}
