package main

import (
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"github.com/spf13/cobra"
	"log"
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

var example = `
  updater CPH2401_11.C.58_0580_202402190800 --model=CPH2401 --region=CN
  updater RMX3820_13.1.0.130_0130_202404010000 --model=RMX3820 --region=IN --mode=1
  updater A127_13.0_0001 --model=A127 --carrier=00000000 --proxy=http://localhost:7890
  updater OPD2413_11.A --model=OPD2413 --region=CN --gray=1
  updater PJX110_11.C --region=CN --reqmode=taste
`

var rootCmd = &cobra.Command{
	Use:     "updater [OTA_VERSION]",
	Short:   "Query OPlus, OPPO and Realme Mobile OS version updates using official API",
	Example: example,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		otaVer := args[0]

		model := getStringFlag(cmd, "model")
		carrier := getStringFlag(cmd, "carrier")
		region := getStringFlag(cmd, "region")
		mode := getIntFlag(cmd, "mode")
		imei := getStringFlag(cmd, "imei")
		guid := getStringFlag(cmd, "guid")
		proxy := getStringFlag(cmd, "proxy")
		gray := getIntFlag(cmd, "gray")
		reqmode := getStringFlag(cmd, "reqmode")

		result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
			OtaVersion: otaVer,
			Region:     region,
			Model:      model,
			NvCarrier:  carrier,
			Mode:       mode,
			IMEI:       imei,
			GUID:       guid,
			Proxy:      proxy,
			Gray:       gray,
			ReqMode:    reqmode,
		})
		if err != nil {
			log.Fatalf("Error in QueryUpdate: %v", err)
		}
		result.PrettyPrint()
	},
}

func init() {
	rootCmd.Flags().StringP("model", "m", "", "Device model (required), e.g., RMX3820, CPH2401")
	rootCmd.Flags().StringP("region", "r", "CN", "Server region: CN (default), EU or IN")
	rootCmd.Flags().StringP("carrier", "c", "", "Carrier ID found in `my_manifest/build.prop` file under the `NV_ID` reference, e.g., 01000100")
	rootCmd.Flags().Int("mode", 0, "Mode: 0 (stable, default) or 1 (testing)")
	rootCmd.Flags().StringP("imei", "i", "", "IMEI, e.g., 864290000000000")
	rootCmd.Flags().StringP("guid", "g", "", "GUID, e.g., 1234567890(64 bit)")
	rootCmd.Flags().StringP("proxy", "p", "", "Proxy server, e.g., type://user:password@host:port")
	rootCmd.Flags().Int("gray", 0, "Gray update server: 0 (default) or 1 (use gray server for CN region)")
	rootCmd.Flags().String("reqmode", "manual", "Request Mode: manual (default), server_auto, client_auto or taste. Do not use taste mode together with gray update mode (--gray=1).")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}