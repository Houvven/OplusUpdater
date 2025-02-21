package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"github.com/tidwall/pretty"
	"testing"
)

func TestQueryUpdater(t *testing.T) {
	result := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX3350_11.C",
	})
	fmt.Println(string(pretty.Color(pretty.Pretty(result.DecryptedBodyBytes), nil)))
}
