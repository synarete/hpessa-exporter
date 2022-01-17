// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"fmt"
	"path/filepath"
	"strings"
)

type SysFS struct {
	PseudoFS
}

const (
	sysDefaultMountPoint = "/sys"
)

func NewSysFS() *SysFS {
	return newSysFS(sysDefaultMountPoint)
}

func newSysFS(prefix string) *SysFS {
	return &SysFS{
		PseudoFS: PseudoFS{
			Prefix: prefix,
		},
	}
}

func (sysfs *SysFS) Sub(subpath string) *PseudoFS {
	return &PseudoFS{
		Prefix: filepath.Join(sysfs.Prefix, subpath),
	}
}

func (sysfs *SysFS) BlockDeviceSub(dev string) *PseudoFS {
	return sysfs.Sub(filepath.Join("block", dev))
}

func (sysfs *SysFS) BlockDeviceQueueSub(dev string) *PseudoFS {
	return sysfs.Sub(filepath.Join("block", dev, "queue"))
}

func (sysfs *SysFS) ListBlockDevices() ([]string, error) {
	return sysfs.ReadDir("block")
}

// IsBlock expects /sys/block/<name>/ dir
func (sysfs *SysFS) IsBlock(name string) (bool, error) {
	return sysfs.IsDir("block", name)
}

// BlockStat parses /sys/block/<device>/stat
// https://www.kernel.org/doc/Documentation/block/stat.txt
func (sysfs *SysFS) BlockStat(dev string) (*BlkdevIOStat, error) {
	dat, err := sysfs.ReadFile("block", dev, "stat")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(dat, "\n")
	if len(lines) < 1 || len(lines) > 2 {
		return nil, fmt.Errorf("wrong number of lines in /sys/block/%s/stat", dev)
	}
	return sysfs.parseBlockStatLine(lines[0])
}

func (sysfs *SysFS) parseBlockStatLine(line string) (*BlkdevIOStat, error) {
	var err error
	ret := &BlkdevIOStat{}
	fields, err := sysfs.SplitFields(line, 11)
	if err != nil {
		return ret, err
	}
	if ret.ReadsIOs, err = sysfs.ParseUint64(fields[0]); err != nil {
		return ret, err
	}
	if ret.ReadsMerged, err = sysfs.ParseUint64(fields[1]); err != nil {
		return ret, err
	}
	if ret.ReadsIOs, err = sysfs.ParseMultUint64(fields[2], SectorSize); err != nil {
		return ret, err
	}
	if ret.ReadTimeMS, err = sysfs.ParseUint64(fields[3]); err != nil {
		return ret, err
	}
	if ret.WritesIOs, err = sysfs.ParseUint64(fields[4]); err != nil {
		return ret, err
	}
	if ret.WritesMerged, err = sysfs.ParseUint64(fields[5]); err != nil {
		return ret, err
	}
	if ret.WritesIOs, err = sysfs.ParseMultUint64(fields[6], SectorSize); err != nil {
		return ret, err
	}
	if ret.WriteTimeMS, err = sysfs.ParseUint64(fields[7]); err != nil {
		return ret, err
	}
	if ret.InFlight, err = sysfs.ParseUint64(fields[8]); err != nil {
		return ret, err
	}
	if ret.IOTimeMS, err = sysfs.ParseUint64(fields[9]); err != nil {
		return ret, err
	}
	if ret.WeightedIOTimeMS, err = sysfs.ParseUint64(fields[10]); err != nil {
		return ret, err
	}
	return ret, nil
}

// BlockQueueStats parses /sys/block/<dev>/queue
func (sysfs *SysFS) BlockQueueStats(dev string) (*BlkdevQueueStats, error) {
	var err error
	ret := &BlkdevQueueStats{}
	pfs := sysfs.BlockDeviceQueueSub(dev)

	if ret.AddRandom, err = pfs.ReadFileAsBool("add_random"); err != nil {
		return ret, err
	}
	if ret.DAX, err = pfs.ReadFileAsBool("dax"); err != nil {
		ret.DAX = false
	}
	if ret.DiscardGranularity, err = pfs.ReadFileAsUInt64("discard_granularity"); err != nil {
		return ret, err
	}
	if ret.DiscardMaxHWBytes, err = pfs.ReadFileAsUInt64("discard_max_hw_bytes"); err != nil {
		return ret, err
	}
	if ret.DiscardMaxBytes, err = pfs.ReadFileAsUInt64("discard_max_bytes"); err != nil {
		return ret, err
	}
	if ret.FUA, err = pfs.ReadFileAsBool("fua"); err != nil {
		return ret, err
	}
	if ret.HWSectorSize, err = pfs.ReadFileAsUInt32("hw_sector_size"); err != nil {
		return ret, err
	}
	if ret.IOPoll, err = pfs.ReadFileAsBool("io_poll"); err != nil {
		return ret, err
	}
	if ret.IOPollDelay, err = pfs.ReadFileAsInt("io_poll_delay"); err != nil {
		return ret, err
	}
	if ret.IOTimeout, err = pfs.ReadFileAsUInt64("io_timeout"); err != nil {
		ret.IOTimeout = 0
	}
	if ret.IOStats, err = pfs.ReadFileAsBool("iostats"); err != nil {
		return ret, err
	}
	if ret.LogicalBlockSize, err = pfs.ReadFileAsUInt64("logical_block_size"); err != nil {
		return ret, err
	}
	if ret.MaxHWSectorsKB, err = pfs.ReadFileAsUInt64("max_hw_sectors_kb"); err != nil {
		return ret, err
	}
	if ret.MaxIntegritySegments, err = pfs.ReadFileAsUInt64("max_integrity_segments"); err != nil {
		return ret, err
	}
	if ret.MaxSectorsKB, err = pfs.ReadFileAsUInt64("max_sectors_kb"); err != nil {
		return ret, err
	}
	if ret.MaxSegments, err = pfs.ReadFileAsUInt64("max_segments"); err != nil {
		return ret, err
	}
	if ret.MaxSegmentSize, err = pfs.ReadFileAsUInt64("max_segment_size"); err != nil {
		return ret, err
	}
	if ret.MinimumIOSize, err = pfs.ReadFileAsUInt64("minimum_io_size"); err != nil {
		return ret, err
	}
	if ret.NoMerges, err = pfs.ReadFileAsUInt32("nomerges"); err != nil {
		return ret, err
	}
	if ret.NRRequests, err = pfs.ReadFileAsUInt64("nr_requests"); err != nil {
		return ret, err
	}
	if ret.OptimalIOSize, err = pfs.ReadFileAsUInt64("optimal_io_size"); err != nil {
		return ret, err
	}
	if ret.PhysicalBlockSize, err = pfs.ReadFileAsUInt64("physical_block_size"); err != nil {
		return ret, err
	}
	if ret.ReadAHeadKB, err = pfs.ReadFileAsUInt64("read_ahead_kb"); err != nil {
		return ret, err
	}
	if ret.Rotational, err = pfs.ReadFileAsBool("rotational"); err != nil {
		return ret, err
	}
	if ret.RQAffinity, err = pfs.ReadFileAsUInt32("rq_affinity"); err != nil {
		return ret, err
	}
	if ret.Scheduler, err = pfs.ReadFile("scheduler"); err != nil {
		return ret, err
	}
	if ret.WriteCache, err = pfs.ReadFile("write_cache"); err != nil {
		return ret, err
	}
	if ret.WriteSameMaxBytes, err = pfs.ReadFileAsUInt64("write_same_max_bytes"); err != nil {
		ret.WriteSameMaxBytes = 0
	}
	if ret.WBTLatUSec, err = pfs.ReadFileAsInt("wbt_lat_usec"); err != nil {
		ret.WBTLatUSec = 0
	}
	if ret.Zoned, err = pfs.ReadFile("zoned"); err != nil {
		ret.Zoned = ""
	}
	if ret.ZoneWriteGranularity, err = pfs.ReadFileAsInt("zone_write_granularity"); err != nil {
		ret.ZoneWriteGranularity = 0
	}
	return ret, nil
}

// BlkdevInfo parses raw information from under /sys/block/<dev>/...
func (sysfs *SysFS) BlkdevInfo(dev string) (*BlkdevInfo, error) {
	var err error
	ret := &BlkdevInfo{Name: dev}
	pfs := sysfs.BlockDeviceSub(dev)

	devstr, err := pfs.ReadFile("dev")
	if err == nil {
		nums := strings.Split(devstr, ":")
		if len(nums) == 2 {
			ret.Major, _ = pfs.ParseUint32(nums[0])
			ret.Minor, _ = pfs.ParseUint32(nums[1])
		}
	}

	if ret.Size, err = pfs.ReadFileAsUInt64("size"); err != nil {
		ret.Size = 0
	}
	if ret.Vendor, err = pfs.ReadFileTrim("device", "vendor"); err != nil {
		ret.Vendor = ""
	}
	if ret.Model, err = pfs.ReadFileTrim("device", "model"); err != nil {
		ret.Model = ""
	}
	if ret.Readonly, err = pfs.ReadFileAsBool("ro"); err != nil {
		ret.Readonly = false
	}
	return ret, nil
}
