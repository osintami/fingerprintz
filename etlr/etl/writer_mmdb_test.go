// Copyright © 2023 OSINTAMI. This is not yours.
package etl

import (
	"net"
	"os"
	"testing"

	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/osintami/fingerprintz/common"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestMMDB(t *testing.T) {
	mmdbFile := "/tmp/test.mmdb"
	os.Remove(mmdbFile)

	_, cidr, _ := net.ParseCIDR("1.2.3.4/32")

	mmdbwriter := NewMMDBWriter()
	assert.Equal(t, "mmdb", mmdbwriter.Type())

	mmdbwriter.Create("test")
	entry := mmdbtype.Map{
		"test": mmdbtype.Map{
			"blacklist": mmdbtype.Map{
				"isBlacklisted": mmdbtype.Bool(true),
			},
		},
	}
	err := mmdbwriter.Insert(cidr, entry)
	assert.Nil(t, err)

	info := &ETLJobInfo{snapshotFile: mmdbFile}
	job := NewMockETLJob(&Toolbox{}, &Source{}, info)

	err = mmdbwriter.Load(job)
	assert.Nil(t, err)

	/// read the mmdb file we just created
	mmdbreader, err := common.NewMaxmindReader("/tmp/test.mmdb")
	assert.NotNil(t, mmdbreader)

	data, err := mmdbreader.Lookup(net.ParseIP("1.2.3.4"))
	assert.Nil(t, err)

	result := gjson.GetBytes(data, "test.blacklist.isBlacklisted")
	assert.True(t, result.Bool())

	os.Remove(mmdbFile)
}
