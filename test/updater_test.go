package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"testing"
)

func TestQueryUpdater(t *testing.T) {
	responseCipher := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX3800_11.C",
	})
	fmt.Println(responseCipher)
}
