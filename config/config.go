package config

type Config struct {
	CarrierID        string // Found in `my_manifest/build.prop` file, under the `NV_ID` reference
	Host             string
	Language         string
	PublicKey        string // Found in `com.oplus.app-features.xml` file
	PublicKeyVersion string // Found in `com.oplus.app-features.xml` file
	Version          string // Found in `com.oplus.app-features.xml` file
}

func GetConfig(country string) Config {
	if country == "EU" {
		return Config{
			CarrierID:        "10100111",
			Host:             "component-ota-eu.allawnos.com",
			Language:         "en-GB",
			PublicKey:        publicKeyEU,
			PublicKeyVersion: "1615897067573",
			Version:          "2",
		}
	}

	if country == "IN" {
		return Config{
			CarrierID:        "00011011",
			Host:             "component-ota-in.allawnos.com",
			Language:         "en-IN",
			PublicKey:        publicKeyIN,
			PublicKeyVersion: "1615896309308",
			Version:          "2",
		}
	}

	if country == "US" {
		return Config{
			CarrierID:        "10100001",
			Host:             "UNKNOWN",
			Language:         "en-US",
			PublicKey:        "UNKNOWN",
			PublicKeyVersion: "UNKNOWN",
			Version:          "2",
		}
	}

	// Default to CN
	return Config{
		CarrierID:        "10010111",
		Host:             "component-ota-cn.allawntech.com",
		Language:         "zh-CN",
		PublicKey:        publicKeyCN,
		PublicKeyVersion: "1615879139745",
		Version:          "2",
	}
}
