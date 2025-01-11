package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"updater/config"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

func (header *UpdateRequestHeaders) makeHeader(c config.Config) (http.Header, error) {
	h := http.Header{
		"androidVersion": {header.AndroidVersion},
		"osVersion":      {header.ColorOSVersion},
		"colorOSVersion": {header.ColorOSVersion},
		"otaVersion":     {header.OtaVersion}, // ro.build.version.ota - my_manifest/build.prop
		"deviceId":       {header.DeviceId},
		"model":          {strings.Split(header.OtaVersion, "_")[0]},
		"language":       {c.Language},
		"nvCarrier":      {c.CarrierID},
		"version":        {c.Version},
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
	marshal, _ := json.Marshal(cipher)
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
	c config.Config,
	proxyStr string,
) (ResponseResult, error) {
	body, err := makeBody(key, iv, cipher)
	if err != nil {
		return ResponseResult{}, err
	}
	header, err := updateHeaders.makeHeader(c)
	if err != nil {
		return ResponseResult{}, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: c.Host, Path: "/update/v3"},
		Header: header,
		Body:   io.NopCloser(bytes.NewBuffer(body)),
	}

	client := &http.Client{}
	parseProxy(client, proxyStr)
	resp, err := client.Do(req)
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
	log.Println(string(cipherJson))
}

func parseProxy(client *http.Client, p string) {
	if p == "" {
		return
	}

	parsedURL, err := url.Parse(p)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		log.Fatalf("Invalid proxy URL: %s", p)
	}

	switch parsedURL.Scheme {
	case "socks":
		var proxyAuth *proxy.Auth
		if parsedURL.User != nil {
			proxyAuth = &proxy.Auth{
				User: parsedURL.User.Username(),
				Password: func() string {
					pass, _ := parsedURL.User.Password()
					return pass
				}(),
			}
		}

		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, proxyAuth, proxy.Direct)
		if err != nil {
			log.Fatalf("Error creating proxy dialer: %v", err)
		}
		client.Transport = &http.Transport{
			Dial: dialer.Dial,
		}

	case "http", "https":
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(parsedURL),
		}

	default:
		log.Fatalf("Unsupported proxy scheme: %s", parsedURL.Scheme)
	}
}
