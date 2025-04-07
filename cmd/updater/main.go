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

var rootCmd = &cobra.Command{
	Use:   "updater",
	Short: " Use Oplus official api to query OPlus,OPPO and Realme Mobile 's OS version update.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//Get the value of the flag
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
			log.Printf("Error in QueryUpdate: %v\n", err)
		}
		result.PrettyPrint()
	},
}

func init() {
	rootCmd.Flags().String("region", "CN", "Server zone: CN (default), EU or IN (optional), e.g., --region=CN")
	rootCmd.Flags().String("model", "", "Device model, e.g., --model=RMX3820")
	rootCmd.Flags().String("carrier", "", "Found in `my_manifest/build.prop` file, under the `NV_ID` reference, e.g., --carrier=01000100")
	rootCmd.Flags().Int("mode", 0, "Mode: 0 (stable, default) or 1 (testing), e.g., --mode=0")
	rootCmd.Flags().String("imei", "", "IMEI, e.g., --imei=86429XXXXXXXX98")
	rootCmd.Flags().StringP("proxy", "p", "", "Proxy server, e.g., --proxy=type://@host:port or --proxy=type://user:password@host:port, support http and socks proxy")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
