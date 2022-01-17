package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

func TestSysfsBlkdev(t *testing.T) {
	procfs := devmon.NewProcFS()
	sysfs := devmon.NewSysFS()

	disks, err := procfs.DiskStats()
	assert.NoError(t, err)
	assert.NotNil(t, disks)

	for _, dio := range disks {
		name := dio.DeviceName
		isblk, _ := sysfs.IsBlock(name)
		if !isblk {
			continue
		}
		stat, err := sysfs.BlockStat(name)
		assert.NoError(t, err)
		assert.NotNil(t, stat)

		qstat, err := sysfs.BlockQueueStats(name)
		assert.NoError(t, err)
		assert.NotNil(t, qstat)

		bdi, err := sysfs.BlkdevInfo(name)
		assert.NoError(t, err)
		assert.NotNil(t, qstat)
		assert.GreaterOrEqual(t, int64(bdi.Size), int64(0))
	}
}
