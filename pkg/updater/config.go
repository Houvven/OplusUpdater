package updater

type Config struct {
	CarrierID        string // Found in `my_manifest/build.prop` file, under the `NV_ID` reference
	Host             string
	Language         string
	PublicKey        string // Found in `com.oplus.app-features.xml` file
	PublicKeyVersion string // Found in `com.oplus.app-features.xml` file
	Version          string // Found in `com.oplus.app-features.xml` file
}

type Region = string

const (
	RegionCn = "CN"
	RegionEu = "EU"
	RegionIn = "IN"
	RegionSg = "SG"
	RegionRu = "RU"
	RegionTr = "TR"
	RegionTh = "TH"
	RegionGl = "GL"
)

const commonHost = "component-ota-sg.allawnos.com"

func GetConfig(region string, gray int) *Config {
	const defaultVersion = "2"

	baseConfigSG := &Config{
		Host:             commonHost,
		PublicKey:        publicKeySG,
		PublicKeyVersion: "1615895993238",
		Version:          defaultVersion,
	}

	regionOverrides := map[string]struct {
		CarrierID string
		Language  string
	}{
		RegionSg: {"01011010", "en-SG"},
		RegionRu: {"00110111", "ru-RU"},
		RegionTr: {"01010001", "tr-TR"},
		RegionTh: {"00111001", "th-TH"},
		RegionGl: {"10100111", "en-US"},
	}

	if override, ok := regionOverrides[region]; ok {
		cfg := *baseConfigSG
		cfg.CarrierID = override.CarrierID
		cfg.Language = override.Language
		return &cfg
	}

	switch region {
	case RegionEu:
		return &Config{
			CarrierID:        "01000100",
			Host:             "component-ota-eu.allawnos.com",
			Language:         "en-GB",
			PublicKey:        publicKeyEU,
			PublicKeyVersion: "1615897067573",
			Version:          defaultVersion,
		}
	case RegionIn:
		return &Config{
			CarrierID:        "00011011",
			Host:             "component-ota-in.allawnos.com",
			Language:         "en-IN",
			PublicKey:        publicKeyIN,
			PublicKeyVersion: "1615896309308",
			Version:          defaultVersion,
		}
	}

	host := "component-ota-cn.allawntech.com"
	if gray == 1 {
		host = "component-ota-gray.coloros.com"
	}

	return &Config{
		CarrierID:        "10010111",
		Host:             host,
		Language:         "zh-CN",
		PublicKey:        publicKeyCN,
		PublicKeyVersion: "1615879139745",
		Version:          defaultVersion,
	}
}