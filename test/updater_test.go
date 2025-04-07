package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"testing"
)

func TestQueryUpdate_CPH2653_EU(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "CPH2653_11.A",
		Region:     updater.RegionEu,
		Model:      "CPH2653EEA",
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}

func TestQueryUpdate_CPH2653_EU_TEST(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "CPH2653_11.A",
		Region:     updater.RegionEu,
		Model:      "CPH2653EEA",
		Mode:       1,
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}

func TestQueryUpdate_RMX5010_CN(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5010_11.A",
		Region:     updater.RegionCn,
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}

func TestQueryUpdate_RMX5011_RU(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5011_11.A",
		Region:     updater.RegionRu,
		Model:      "RMX5011RU",
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}

func TestQueryUpdate_RMX5011_TR(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5011_11.A",
		Region:     updater.RegionTr,
		Model:      "RMX5011TR",
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}

func TestQueryUpdate_PHP110_CN(t *testing.T) {
	result, err := updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "PHP110_11.F",
		Region:     updater.RegionCn,
	})
	if err != nil {
		fmt.Println(err)
	}
	result.PrettyPrint()
}
