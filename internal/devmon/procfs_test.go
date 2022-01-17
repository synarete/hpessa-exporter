// SPDX-License-Identifier: Apache-2.0
package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

func TestProcfsLoadAvg(t *testing.T) {
	procfs := devmon.NewProcFS()
	assert.NotNil(t, procfs)
	loadavg, err := procfs.LoadAvg()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, loadavg.Load1, 0.0)
}

func TestProcfsDiskStats(t *testing.T) {
	procfs := devmon.NewProcFS()
	assert.NotNil(t, procfs)
	stats, err := procfs.DiskStats()
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Greater(t, len(stats), 0)
	var readIOs uint64
	for _, ios := range stats {
		readIOs += ios.ReadsIOs
	}
	assert.Greater(t, readIOs, uint64(0))
}
