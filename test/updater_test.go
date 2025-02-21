package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"github.com/tidwall/pretty"
	"testing"
)

func TestQueryUpdater(t *testing.T) {
	result := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "CPH2653_11.A",
		Region:     updater.RegionEu,
		Model:      "CPH2653EEA",
	})

	fmt.Printf("Status: %d\n", result.ResponseCode)
	if result.ResponseCode != 200 {
		fmt.Printf("Error: %s\n", result.ErrMsg)
	}
	fmt.Println(string(pretty.Color(pretty.Pretty(result.DecryptedBodyBytes), nil)))
}
