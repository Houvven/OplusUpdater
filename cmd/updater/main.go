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
		proxy := getStringFlag(cmd, "proxy")

		result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
			OtaVersion: otaVer,
			Region:     region,
			Model:      model,
			NvCarrier:  carrier,
			Mode:       mode,
			IMEI:       imei,
			Proxy:      proxy,
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
	rootCmd.Flags().StringP("proxy", "p", "", "Proxy server, e.g., type://user:password@host:port")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
