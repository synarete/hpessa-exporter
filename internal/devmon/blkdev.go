// SPDX-License-Identifier: Apache-2.0
package devmon

type BlkdevMap struct {
	IDs     map[string]BlkdevID
	IOStats map[string]*BlkdevIOStat
	QStats  map[string]*BlkdevQueueStats
}

func DiscoverBlkdevStats(procfs *ProcFS, sysfs *SysFS) (*BlkdevMap, error) {
	bdm := &BlkdevMap{
		IDs:     map[string]BlkdevID{},
		IOStats: map[string]*BlkdevIOStat{},
		QStats:  map[string]*BlkdevQueueStats{},
	}
	disks, err := procfs.DiskStats()
	if err != nil {
		return bdm, err
	}
	for _, dio := range disks {
		name := dio.DeviceName
		bdm.IDs[name] = dio.BlkdevID

		isblk, _ := sysfs.IsBlock(name)
		if !isblk {
			continue
		}

		stat, err := sysfs.BlockStat(name)
		if err == nil {
			bdm.IOStats[name] = stat
		}

		qstat, err := sysfs.BlockQueueStats(name)
		if err == nil {
			bdm.QStats[name] = qstat
		}
	}
	return bdm, nil
}

func DiscoverBlkdevInfo(procfs *ProcFS, sysfs *SysFS) ([]BlkdevInfo, error) {
	ret := []BlkdevInfo{}
	disks, err := procfs.DiskStats()
	if err != nil {
		return ret, err
	}
	for _, dio := range disks {
		name := dio.DeviceName
		isblk, _ := sysfs.IsBlock(name)
		if !isblk {
			continue
		}
		bdi, err := sysfs.BlkdevInfo(name)
		if err == nil {
			ret = append(ret, *bdi)
		}
	}
	return ret, nil
}

func DescoveBlockDevicesIO(procfs *ProcFS, sysfs *SysFS) ([]BlkdevIOInfo, error) {
	ret := []BlkdevIOInfo{}
	disks, err := procfs.DiskStats()
	if err != nil {
		return ret, err
	}
	for _, dio := range disks {
		name := dio.DeviceName
		if isblk, _ := sysfs.IsBlock(name); isblk {
			ret = append(ret, dio)
		}
	}
	return ret, nil
}
