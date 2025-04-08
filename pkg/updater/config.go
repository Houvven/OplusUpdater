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
)

const commonHost = "component-ota-sg.allawnos.com"

func GetConfig(region string) *Config {

	if region == RegionEu {
		return &Config{
			CarrierID:        "01000100",
			Host:             "component-ota-eu.allawnos.com",
			Language:         "en-GB",
			PublicKey:        publicKeyEU,
			PublicKeyVersion: "1615897067573",
			Version:          "2",
		}
	}

	if region == RegionIn {
		return &Config{
			CarrierID:        "00011011",
			Host:             "component-ota-in.allawnos.com",
			Language:         "en-IN",
			PublicKey:        publicKeyIN,
			PublicKeyVersion: "1615896309308",
			Version:          "2",
		}
	}

	if region == RegionSg || region == RegionRu || region == RegionTr || region == RegionTh {
		c := &Config{
			Host:             commonHost,
			PublicKey:        publicKeySG,
			PublicKeyVersion: "1615895993238",
			Version:          "2",
		}
		if region == RegionRu {
			c.CarrierID = "00110111"
			c.Language = "ru-RU"
		} else if region == RegionTr {
			c.CarrierID = "01010001"
			c.Language = "tr-TR"
		} else if region == RegionTh {
			c.CarrierID = "00111001"
			c.Language = "th-TH"
		} else {
			c.CarrierID = "01011010"
			c.Language = "en-SG"
		}
		return c
	}

	// Default to CN
	return &Config{
		CarrierID:        "10010111",
		Host:             "component-ota-cn.allawntech.com",
		Language:         "zh-CN",
		PublicKey:        publicKeyCN,
		PublicKeyVersion: "1615879139745",
		Version:          "2",
	}
}
