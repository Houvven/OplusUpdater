package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"testing"
)

func TestCalculateIMEICheckDigit(t *testing.T) {
	imei := "86429007457802"
	digit, err := updater.CalculateIMEICheckDigit(imei)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("IMEI: %s, Check Digit: %s\n", imei, digit)
}
