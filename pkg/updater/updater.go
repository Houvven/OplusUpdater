package updater

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Attribute struct {
	Zone       string
	Mode       int
	OtaVer     string
	AndroidVer string
	ColorOSVer string
	Transport  *http.Transport
}

func (attr *Attribute) postProcessing() {
	// maybe only support at realme.
	// i dont known oppo and oplus 's version format.
	if len(strings.Split(attr.OtaVer, "_")) < 3 || len(strings.Split(attr.OtaVer, ".")) < 3 {
		attr.OtaVer += ".00_0000_000000000000"
	}
	if attr.Zone == "" {
		attr.Zone = "CN"
	}
	if attr.AndroidVer == "" {
		attr.AndroidVer = "nil"
	}
	if attr.ColorOSVer == "" {
		attr.ColorOSVer = "nil"
	}
	if attr.Mode == 0 {
		attr.Mode = 0
	}
	if attr.Transport == nil {
		attr.Transport = &http.Transport{}
	}
}

type UpdateResponseCipher struct {
	Components []struct {
		ComponentId      string `json:"componentId"`
		ComponentName    string `json:"componentName"`
		ComponentVersion string `json:"componentVersion"`

		ComponentPackets struct {
			Size      string `json:"size"`
			ManualUrl string `json:"manualUrl"`
			Id        string `json:"id"`
			Url       string `json:"url"`
			Md5       string `json:"md5"`
		} `json:"componentPackets"`
	} `json:"components"`

	SecurityPatch   string `json:"securityPatch"`
	RealVersionName string `json:"realVersionName"`

	Description struct {
		PanelUrl   string `json:"panelUrl"`
		Url        string `json:"url"`
		FirstTitle string `json:"firstTitle"`
	} `json:"description"`

	RealAndroidVersion  string `json:"realAndroidVersion"`
	RealOsVersion       string `json:"realOsVersion"`
	SecurityPatchVendor string `json:"securityPatchVendor"`
	RealOtaVersion      string `json:"realOtaVersion"`
	VersionTypeH5       string `json:"versionTypeH5"`
	Status              string `json:"status"`
}

// func QueryUpdater(otaVer, androidVer, colorOsVer, zone string, mode int, transport http.Transport) (*UpdateResponseCipher, error) {

func QueryUpdater(attr Attribute) (*UpdateResponseCipher, error) {
	rawBytes, err := QueryUpdaterRawBytes(attr)
	if err != nil {
		return nil, err
	}
	var cipher UpdateResponseCipher
	if err := json.Unmarshal(rawBytes, &cipher); err != nil {
		return nil, err
	}
	return &cipher, nil
}

func QueryUpdaterRawBytes(attr Attribute) ([]byte, error) {
	attr.postProcessing()
	deviceId := GetDefaultDeviceId()

	c := GetConfig(attr.Zone)
	key, err := RandomKey()
	if err != nil {
		return nil, err
	}
	iv, err := RandomIv()
	if err != nil {
		return nil, err
	}
	protectedKey, err := GenerateProtectedKey(key, []byte(c.PublicKey))
	if err != nil {
		return nil, err
	}

	headers := UpdateRequestHeaders{
		AndroidVersion: attr.AndroidVer, // or Android13
		ColorOSVersion: attr.ColorOSVer, // or ColorOS13.1.0
		OtaVersion:     attr.OtaVer,
		ProtectedKey: map[string]CryptoConfig{
			"SCENE_1": {
				ProtectedKey:       protectedKey,
				Version:            GenerateProtectedVersion(),
				NegotiationVersion: c.PublicKeyVersion,
			},
		},
	}
	headers.SetHashedDeviceId(deviceId)
	cipher := NewUpdateRequestCipher(attr.Mode, deviceId)

	reqHeaders, err := headers.CreateRequestHeader(c)
	if err != nil {
		return nil, err
	}
	reqBody, err := cipher.CreateRequestBody(key, iv)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: c.Host, Path: "/update/v3"},
		Header: reqHeaders,
		Body:   io.NopCloser(bytes.NewBuffer(reqBody)),
	}
	client := &http.Client{
		Transport: attr.Transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {}(req.Body)

	var result ResponseResult
	if json.NewDecoder(resp.Body).Decode(&result) != nil {
		return nil, err
	}

	return decryptUpdateResponse(&result, key)
}

func decryptUpdateResponse(r *ResponseResult, key []byte) ([]byte, error) {
	if r.ResponseCode != 200 {
		return nil, fmt.Errorf("respnse code: %d, message: %s", r.ResponseCode, r.ErrMsg)
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(r.Body.(string)), &m); err != nil {
		return nil, err
	}

	iv, err := base64.StdEncoding.DecodeString(m["iv"].(string))
	if err != nil {
		return nil, err
	}
	cipherBytes := crypto.FromBase64String(m["cipher"].(string)).
		Aes().CTR().NoPadding().
		WithKey(key).WithIv(iv).
		Decrypt().ToBytes()

	return cipherBytes, nil
}
