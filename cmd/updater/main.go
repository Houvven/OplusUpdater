package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Houvven/OplusUpdater/pkg/updater"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

func getStringFlag(cmd *cobra.Command, flagName string) string {
	flag, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatalf("Error in %s: %v", flagName, err)
	}
	return flag
}

func getIntFlag(cmd *cobra.Command, flagName string) int {
	flag, err := cmd.Flags().GetInt(flagName)
	if err != nil {
		log.Fatalf("Error in %s: %v", flagName, err)
	}
	return flag
}

var rootCmd = &cobra.Command{
	Use:   "oplus-updater",
	Short: " Use Oplus official api to query OPlus,OPPO and Realme Mobile 's OS version update.",
	Run: func(cmd *cobra.Command, args []string) {
		//Get the value of the flag
		model := getStringFlag(cmd, "model")
		otaVer := getStringFlag(cmd, "ota-version")
		zone := getStringFlag(cmd, "zone")
		mode := getIntFlag(cmd, "mode")
		proxy := getStringFlag(cmd, "proxy")

		result := updater.QueryUpdate(&updater.QueryUpdateArgs{
			Model:      model,
			OtaVersion: otaVer,
			Region:     zone,
			Mode:       mode,
			Proxy:      proxy,
		})
		fmt.Printf("Status: %d\n", result.ResponseCode)
		if result.ResponseCode != 200 {
			fmt.Printf("Error: %s\n", result.ErrMsg)
		}
		fmt.Println(string(pretty.Color(pretty.Pretty(result.DecryptedBodyBytes), nil)))
	},
}

func init() {
	otaVerBytes, _ := exec.Command("getprop", "ro.build.version.ota").Output()
	otaVer := strings.TrimSpace(string(otaVerBytes))

	rootCmd.Flags().StringP("ota-version", "o", otaVer, "OTA version (required), e.g., --ota-version=RMX3820_11.A.00_0000_000000000000")
	rootCmd.Flags().String("region", "CN", "Server zone: CN (default), EU or IN (optional), e.g., --zone=CN")
	rootCmd.Flags().IntP("mode", "m", 0, "Mode: 0 (stable, default) or 1 (testing), e.g., --mode=0")
	rootCmd.Flags().String("model", "", "Device model, e.g., --model=RMX3820")
	rootCmd.Flags().StringP("proxy", "p", "", "Proxy server, e.g., --proxy=type://@host:port or --proxy=type://user:password@host:port, support http and socks proxy")

	if err := rootCmd.MarkFlagRequired("ota-version"); err != nil {
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
