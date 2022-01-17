// SPDX-License-Identifier: Apache-2.0
package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

func TestDiscoverBlockDevices(t *testing.T) {
	bdm, err := devmon.DiscoverBlkdevStats(devmon.NewProcFS(), devmon.NewSysFS())
	assert.NoError(t, err)
	assert.NotNil(t, bdm)
	assert.Greater(t, len(bdm.IDs), 0)
	assert.Greater(t, len(bdm.IOStats), 0)
	assert.Greater(t, len(bdm.QStats), 0)
	for _, dev := range bdm.IDs {
		assert.Greater(t, len(dev.DeviceName), 0)
		assert.Greater(t, dev.MajorNumber, uint32(0))
	}
}

func TestDiscoverBlkdevInfo(t *testing.T) {
	bdi, err := devmon.DiscoverBlkdevInfo(devmon.NewProcFS(), devmon.NewSysFS())
	assert.NoError(t, err)
	assert.Greater(t, len(bdi), 0)
	for _, di := range bdi {
		assert.Greater(t, len(di.Name), 0)
		assert.Greater(t, di.Major, uint32(0))
	}
}
