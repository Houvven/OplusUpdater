package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func (header *UpdateRequestHeaders) SetDeviceId(id string) {
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

func DefaultUpdateRequestCipher() (RequestCipher, error) {
	var params RequestCipher
	err := json.Unmarshal([]byte(DefaultRequestBodyJsonStr), &params)
	return params, err
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
