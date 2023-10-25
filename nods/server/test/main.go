package main

import (
	"encoding/json"
	"net"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
)

func main() {
	CreateMockFastDBs()
	CreateIpsumMMDB()
	CreateUHBMMDB()
	CreateMaxmindDB()
}

func CreateIpsumMMDB() {
	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	mmdbwriter := common.NewMaxmindWriter("ipsum")
	entry := mmdbtype.Map{
		"ipsum": mmdbtype.Map{
			"blacklist": mmdbtype.Map{
				"isBlacklisted":  mmdbtype.Bool(true),
				"lastReportedAt": mmdbtype.String("1942-01-01 00:00:00"),
			},
		},
	}
	mmdbwriter.Insert(cidr, entry)
	mmdbwriter.Close("ipsum.mmdb")
}

func CreateUHBMMDB() {
	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	mmdbwriter := common.NewMaxmindWriter("uhb")
	entry := mmdbtype.Map{
		"uhb": mmdbtype.Map{
			"blacklist": mmdbtype.Map{
				"isBlacklisted":  mmdbtype.Bool(true),
				"lastReportedAt": mmdbtype.String("1942-01-01 00:00:00"),
			},
		},
	}
	mmdbwriter.Insert(cidr, entry)
	mmdbwriter.Close("uhb.mmdb")
}

func CreateMockFastDBs() {
	type FakeFilterInfo struct {
		Domain         string
		IsFake         bool
		LastReportedAt string
	}
	type FakeFilterRow struct {
		Key    string         `json:"Key"`
		Result FakeFilterInfo `json:"Result"`
	}

	cache := common.NewPersistentCache("fakefilter.fast")

	row := &FakeFilterRow{}
	row.Key = "nope.com"
	row.Result.Domain = "nope.com"
	row.Result.IsFake = true
	row.Result.LastReportedAt = "1942-01-01 00:00:00"
	raw, _ := json.Marshal(row)
	cache.Set(row.Key, raw, -1)
	cache.Persist()
}

func CreateMaxmindDB() {
	mmdbwriter := common.NewMaxmindWriter("maxmind")

	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	entry := mmdbtype.Map{
		"continent": mmdbtype.Map{
			"code":       mmdbtype.String("AS"),
			"geoname_id": mmdbtype.Int32(6255147),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("Asia"),
			},
		},
		"country": mmdbtype.Map{
			"geoname_id": mmdbtype.Int32(1835841),
			"iso_code":   mmdbtype.String("KR"),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("South Korea"),
			},
		},
		"location": mmdbtype.Map{
			"accuracy_radius": mmdbtype.Uint16(200),
			"latitude":        mmdbtype.Float64(37.511200),
			"longitude":       mmdbtype.Float64(126.974100),
			"time_zone":       mmdbtype.String("Asia/Seoul"),
		},
		"registered_country": mmdbtype.Map{
			"geoname_id": mmdbtype.Int32(1835841),
			"iso_code":   mmdbtype.String("KR"),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("South Korea"),
			},
		},
	}

	mmdbwriter.Insert(cidr, entry)

	_, cidr, _ = net.ParseCIDR("4.3.2.1/32")

	entry = mmdbtype.Map{
		"continent": mmdbtype.Map{
			"code":       mmdbtype.String("NA"),
			"geoname_id": mmdbtype.Int32(6255147),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("North America"),
			},
		},
		"country": mmdbtype.Map{
			"geoname_id": mmdbtype.Int32(1835841),
			"iso_code":   mmdbtype.String("US"),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("United States"),
			},
		},
		"location": mmdbtype.Map{
			"accuracy_radius": mmdbtype.Uint16(200),
			"latitude":        mmdbtype.Float64(30.633263),
			"longitude":       mmdbtype.Float64(-97.677986),
			"time_zone":       mmdbtype.String("Chicago/CST"),
		},
		"registered_country": mmdbtype.Map{
			"geoname_id": mmdbtype.Int32(1835841),
			"iso_code":   mmdbtype.String("US"),
			"names": mmdbtype.Map{
				"en": mmdbtype.String("United States"),
			},
		},
	}

	mmdbwriter.Insert(cidr, entry)

	mmdbwriter.Close("maxmind.mmdb")
}
