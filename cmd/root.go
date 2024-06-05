package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
	"updater/api"
	"updater/util"
)

func handlerError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "oplus-updater",
	Short: " Use Oplus official api to query OPlus,OPPO and Realme Mobile 's OS version update.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the value of the flag
		otaVer, err := cmd.Flags().GetString("ota-version")
		handlerError(err)
		androidVer, err := cmd.Flags().GetString("android-version")
		handlerError(err)
		colorOSVer, err := cmd.Flags().GetString("colorOS-version")
		handlerError(err)

		do(otaVer, androidVer, colorOSVer)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	otaVerBytes, err := exec.Command("getprop", "ro.build.version.ota").Output()

	rootCmd.Flags().StringP("ota-version", "o", strings.TrimSpace(string(otaVerBytes)), "OTA version (required), e.g., --ota-version=RMX3820_11.A.00_0000_000000000000")
	rootCmd.Flags().StringP("android-version", "a", "nil", "Android version (optional), e.g., --android-version=Android14")
	rootCmd.Flags().StringP("colorOS-version", "c", "nil", "ColorOS version (optional), e.g., --colorOS-version=ColorOS14.1.0")

	if err != nil {
		if err := rootCmd.MarkFlagRequired("ota-version"); err != nil {
			os.Exit(1)
		}
	}
}

func do(otaVer, androidVer, colorOSVer string) {
	var handlerErr = func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	key, err := util.RandomKey()
	handlerErr(err)

	iv, err := util.RandomIv()
	handlerErr(err)

	protectedKey, err := util.GenerateProtectedKey(key, []byte(api.PublicKey))
	handlerErr(err)

	updateHeaders := api.UpdateRequestHeaders{
		AndroidVersion: androidVer, // or Android13
		ColorOSVersion: colorOSVer, // or ColorOS13.1.0
		OtaVersion:     otaVer,
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
	handlerErr(err)
	cipher.DeviceId = deviceId

	result, err := api.RequestUpdate(key, iv, updateHeaders, cipher)
	handlerErr(err)

	api.UpdateResponseParse(result, key)
}
