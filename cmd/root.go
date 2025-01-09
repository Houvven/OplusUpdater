package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"updater/api"
	"updater/config"
	"updater/util"

	"github.com/spf13/cobra"
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
		pattern := `^[A-Z]{4}\d{4}_\d{2}\.[A-Z]$`
		re := regexp.MustCompile(pattern)
		if re.MatchString(otaVer) {
			otaVer += ".00_0000_000000000000"
		}

		handlerError(err)
		androidVer, err := cmd.Flags().GetString("android-version")
		handlerError(err)
		colorOSVer, err := cmd.Flags().GetString("colorOS-version")
		handlerError(err)
		zone, err := cmd.Flags().GetString("zone")
		handlerError(err)
		mode, err := cmd.Flags().GetString("mode")
		handlerError(err)

		do(zone, mode, otaVer, androidVer, colorOSVer)
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
	rootCmd.Flags().StringP("zone", "z", "CN", "Server zone: CN (default), EU or IN (optional), e.g., --zone=CN")
	rootCmd.Flags().StringP("mode", "m", "0", "Mode: 0 (stable, default) or 1 (testing), e.g., --mode=0")

	if err != nil {
		if err := rootCmd.MarkFlagRequired("ota-version"); err != nil {
			os.Exit(1)
		}
	}
}

func do(zone, mode, otaVer, androidVer, colorOSVer string) {
	var handlerErr = func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	c := config.GetConfig(zone)

	key, err := util.RandomKey()
	handlerErr(err)

	iv, err := util.RandomIv()
	handlerErr(err)

	protectedKey, err := util.GenerateProtectedKey(key, []byte(c.PublicKey))
	handlerErr(err)

	updateHeaders := api.UpdateRequestHeaders{
		AndroidVersion: androidVer, // or Android13
		ColorOSVersion: colorOSVer, // or ColorOS13.1.0
		OtaVersion:     otaVer,
		ProtectedKey: map[string]api.CryptoConfig{
			"SCENE_1": {
				ProtectedKey:       protectedKey,
				Version:            util.GenerateProtectedVersion(),
				NegotiationVersion: c.PublicKeyVersion,
			},
		},
	}

	deviceId := util.GetDefaultDeviceId()
	updateHeaders.SetDeviceId(deviceId)

	cipher, err := api.DefaultUpdateRequestCipher(mode)
	handlerErr(err)
	cipher.DeviceId = deviceId

	result, err := api.RequestUpdate(key, iv, updateHeaders, cipher, c)
	handlerErr(err)

	api.UpdateResponseParse(result, key)
}
