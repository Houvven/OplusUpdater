package updater

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"net/http"
	"strings"
)

type ResponseResult struct {
	ResponseCode int    `json:"responseCode"`
	ErrMsg       string `json:"errMsg"`
	Body         any    `json:"body"`
}

type CryptoConfig struct {
	ProtectedKey       string `json:"protectedKey"`
	Version            string `json:"version"`
	NegotiationVersion string `json:"negotiationVersion"`
}

type UpdateRequestHeaders struct {
	AndroidVersion string                  `json:"androidVersion"`
	ColorOSVersion string                  `json:"colorOSVersion"`
	OtaVersion     string                  `json:"otaVersion"`
	DeviceId       string                  `json:"deviceId"`
	ProtectedKey   map[string]CryptoConfig `json:"protectedKey"`
}

func (header *UpdateRequestHeaders) CreateRequestHeader(c Config) (http.Header, error) {
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

func (header *UpdateRequestHeaders) SetHashedDeviceId(id string) {
	hash := sha256.Sum256([]byte(id))
	header.DeviceId = strings.ToUpper(hex.EncodeToString(hash[:]))
}

type RequestCipher struct {
	Components []struct {
		ComponentName    string `json:"componentName"`
		ComponentVersion string `json:"componentVersion"`
	} `json:"components"`

	Mode                int    `json:"mode"`
	Time                int64  `json:"time"`
	IsRooted            string `json:"isRooted"`
	Type                string `json:"type"`
	RegistrationId      string `json:"registrationId"`
	SecurityPatch       string `json:"securityPatch"`
	SecurityPatchVendor string `json:"securityPatchVendor"`
	StrategyVersion     string `json:"strategyVersion"`

	Cota struct {
		CotaVersion     string `json:"cotaVersion"`
		CotaVersionName string `json:"cotaVersionName"`
		BuildType       string `json:"buildType"`
	} `json:"cota"`

	Opex struct {
		Check bool `json:"check"`
	} `json:"opex"`

	Sota struct {
		SotaProtocolVersion string `json:"sotaProtocolVersion"`
		SotaVersion         string `json:"sotaVersion"`
		OtaUpdateTime       int    `json:"otaUpdateTime"`
		UpdateViaReboot     int    `json:"updateViaReboot"`
	} `json:"sota"`

	DeviceId      string `json:"deviceId"`
	Duid          string `json:"duid"`
	H5LinkVersion int    `json:"h5LinkVersion"`
}

func NewUpdateRequestCipher(mode int, deviceId string) *RequestCipher {
	return &RequestCipher{
		Mode:     mode,
		IsRooted: "0",
		DeviceId: deviceId,
	}
}

func (cipher *RequestCipher) CreateRequestBody(key, iv []byte) ([]byte, error) {
	marshal, err := json.Marshal(cipher)
	if err != nil {
		return nil, err
	}

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
