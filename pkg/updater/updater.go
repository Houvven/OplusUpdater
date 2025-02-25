package updater

import (
	"encoding/base64"
	"encoding/json"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/go-resty/resty/v2"
	"net/url"
	"strings"
	"time"
)

type QueryUpdateArgs struct {
	OtaVersion string
	Region     string
	Model      string
	NvCarrier  string
	Mode       int
	IMEI       string
	Proxy      string
}

func (args *QueryUpdateArgs) post() {
	if len(strings.Split(args.OtaVersion, "_")) < 3 || len(strings.Split(args.OtaVersion, ".")) < 3 {
		args.OtaVersion += ".00_0000_000000000000"
	}
	if r := strings.TrimSpace(args.Region); len(r) == 0 {
		args.Region = RegionCn
	}
	if m := strings.TrimSpace(args.Model); len(m) == 0 {
		args.Model = strings.Split(args.OtaVersion, "_")[0]
	}
}

func QueryUpdate(args *QueryUpdateArgs) (*ResponseResult, error) {
	args.post()

	config := GetConfig(args.Region)
	if args.NvCarrier == "" {
		args.NvCarrier = config.CarrierID
	}
	iv, err := RandomIv()
	if err != nil {
		return nil, err
	}
	key, err := RandomKey()
	if err != nil {
		return nil, err
	}
	protectedKey, err := GenerateProtectedKey(key, []byte(config.PublicKey))
	if err != nil {
		return nil, err
	}

	var deviceId string
	if len(strings.TrimSpace(args.IMEI)) == 0 {
		deviceId = GenerateDefaultDeviceId()
	} else {
		deviceId = GenerateDeviceId(args.IMEI)
	}

	requestUrl := url.URL{Host: config.Host, Scheme: "https", Path: "/update/v5"}
	requestHeaders := map[string]string{
		"language":       config.Language,
		"androidVersion": "unknown",
		"colorOSVersion": "unknown",
		"otaVersion":     args.OtaVersion,
		"model":          args.Model,
		"nvCarrier":      args.NvCarrier,
		"version":        config.Version,
		"deviceId":       deviceId,
		"Content-Type":   "application/json; charset=utf-8",
	}
	pkm := map[string]CryptoConfig{
		"SCENE_1": {
			ProtectedKey:       protectedKey,
			Version:            GenerateProtectedVersion(),
			NegotiationVersion: config.PublicKeyVersion,
		},
	}
	if pk, err := json.Marshal(pkm); err == nil {
		requestHeaders["protectedKey"] = string(pk)
	} else {
		return nil, err
	}

	var requestBody string
	if r, err := json.Marshal(map[string]any{
		"mode":     args.Mode,
		"time":     time.Now().UnixMilli(),
		"isRooted": "0",
		"isLocked": true,
		"type":     "1",
		"deviceId": deviceId,
	}); err == nil {
		bytes, err := json.Marshal(RequestBody{
			Cipher: crypto.FromBytes(r).
				Aes().CTR().NoPadding().
				WithKey(key).WithIv(iv).
				Encrypt().
				ToBase64String(),
			Iv: base64.StdEncoding.EncodeToString(iv),
		})
		if err != nil {
			return nil, err
		} else {
			requestBody = string(bytes)
		}
	} else {
		return nil, err
	}

	client := resty.New()
	if p := strings.TrimSpace(args.Proxy); len(p) > 0 {
		client.SetProxy(p)
	}
	response, err := client.R().
		SetHeaders(requestHeaders).
		SetBody(map[string]string{"params": requestBody}).
		Post(requestUrl.String())

	if err != nil {
		return nil, err
	}

	var responseResult *ResponseResult
	if json.Unmarshal(response.Body(), &responseResult) != nil {
		return nil, err
	}

	if err := responseResult.DecryptBody(key); err != nil {
		return nil, err
	}

	return responseResult, nil
}
