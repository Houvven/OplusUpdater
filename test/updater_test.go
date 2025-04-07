package test

import (
	"fmt"
	"github.com/Houvven/OplusUpdater/pkg/updater"
	"testing"
)

func ProcessResult(result *updater.ResponseResult, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	result.PrettyPrint()
}

func TestQueryUpdate_CPH2653_EU(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "CPH2653_11.A",
		Region:     updater.RegionEu,
		Model:      "CPH2653EEA",
	}))
}

func TestQueryUpdate_CPH2653_EU_TEST(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "CPH2653_11.A",
		Region:     updater.RegionEu,
		Model:      "CPH2653EEA",
		Mode:       1,
	}))
}

func TestQueryUpdate_RMX5010_CN(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5010_11.A",
		Region:     updater.RegionCn,
	}))
}

func TestQueryUpdate_RMX5011_RU(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5011_11.A",
		Region:     updater.RegionRu,
		Model:      "RMX5011RU",
	}))
}

func TestQueryUpdate_RMX5011_TR(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5011_11.A",
		Region:     updater.RegionTr,
		Model:      "RMX5011TR",
	}))
}

func TestQueryUpdate_RMX5011_TH(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "RMX5011_11.A",
		Region:     updater.RegionTh,
		Model:      "RMX5011",
	}))
}

func TestQueryUpdate_PHP110_CN(t *testing.T) {
	ProcessResult(updater.QueryUpdate(&updater.QueryUpdateArgs{
		OtaVersion: "PHP110_11.F",
		Region:     updater.RegionCn,
	}))
}
