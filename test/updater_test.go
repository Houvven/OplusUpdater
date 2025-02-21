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
	// string to bytes

	var bodyBytes []byte
	switch v := result.Body.(type) {
	case string:
		bodyBytes = []byte(v)
	case []byte:
		bodyBytes = v
	default:
		t.Fatalf("unexpected type for result.Body: %T", result.Body)
	}

	fmt.Println(string(pretty.Color(pretty.Pretty(bodyBytes), nil)))
}
