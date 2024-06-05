package main

import (
	"log"
	"updater/api"
	"updater/util"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	key, err := util.RandomKey()
	handleError(err)

	iv, err := util.RandomIv()
	handleError(err)

	protectedKey, err := util.GenerateProtectedKey(key, []byte(api.PublicKey))
	handleError(err)

	updateHeaders := api.UpdateRequestHeaders{
		AndroidVersion: "null", // or Android13
		ColorOSVersion: "null", // or ColorOS13.1.0
		OtaVersion:     "RMX3820_11.A.00_0000_000000000000",
		ProtectedKey: map[string]api.CryptoConfig{
			"SCENE_1": {
				ProtectedKey:       protectedKey,
				Version:            util.GenerateProtectedVersion(),
				NegotiationVersion: "1615879139745",
			},
		},
	}

	deviceId := util.GetDefaultDeviceId()
	updateHeaders.SetDeviceId(deviceId)

	cipher, err := api.DefaultUpdateRequestCipher()
	handleError(err)
	cipher.DeviceId = deviceId

	result, err := api.RequestUpdate(key, iv, updateHeaders, cipher)
	handleError(err)

	api.UpdateResponseParse(result, key)
}
