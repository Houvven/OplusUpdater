package test

import (
	"encoding/json"
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"github.com/tidwall/pretty"
	"testing"
)

func TestQueryUpdater(t *testing.T) {
	responseCipher, err := updater.QueryUpdater(&updater.Attribute{
		OtaVer: "RMX3800_11.C",
	})
	if err != nil {
		t.Error(err)
	}
	cipherJson, err := json.MarshalIndent(responseCipher, "", "  ")
	cipherJson = pretty.Color(cipherJson, nil)
	fmt.Println(string(cipherJson))
}
