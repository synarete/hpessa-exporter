// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"fmt"
	"strings"
)

type ProcFS struct {
	PseudoFS
}

const (
	procDefaultMountPoint = "/proc"
)

func NewProcFS() *ProcFS {
	return newProcFS(procDefaultMountPoint)
}

func newProcFS(prefix string) *ProcFS {
	return &ProcFS{
		PseudoFS: PseudoFS{
			Prefix: prefix,
		},
	}
}

func (procfs *ProcFS) LoadAvg() (*SysLoadAvg, error) {
	vals, err := procfs.readLoadAvg()
	if err != nil {
		return nil, err
	}

	return procfs.parseLoadAvg(vals)
}

func (procfs *ProcFS) readLoadAvg() ([]string, error) {
	vals, err := procfs.ReadFileFields("loadavg")
	if err != nil {
		return nil, err
	}

	if len(vals) < 3 {
		return nil, fmt.Errorf("malformed loadavg: %q", vals)
	}

	return vals, nil
}

func (procfs *ProcFS) parseLoadAvg(vals []string) (*SysLoadAvg, error) {
	if len(vals) < 3 {
		return nil, fmt.Errorf("malformed loadavg: %q", vals)
	}

	load1, err := procfs.ParseFloat(vals[0])
	if err != nil {
		return nil, err
	}

	load5, err := procfs.ParseFloat(vals[1])
	if err != nil {
		return nil, err
	}

	load15, err := procfs.ParseFloat(vals[2])
	if err != nil {
		return nil, err
	}

	return &SysLoadAvg{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

// DiskStats converts "/proc/diskstats" info into BlkdevIOStat representation.
// See: https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
func (procfs *ProcFS) DiskStats() ([]BlkdevIOInfo, error) {
	ret := []BlkdevIOInfo{}
	dat, err := procfs.ReadFile("diskstats")
	if err != nil {
		return ret, err
	}
	for _, line := range strings.Split(dat, "\n") {
		ln := strings.TrimSpace(line)
		if len(ln) == 0 {
			continue
		}
		ios, err := procfs.parseDiskStatsLine(ln)
		if err != nil {
			return ret, err
		}
		if ios == nil {
			continue
		}
		ret = append(ret, *ios)
	}

	return ret, nil
}

func (procfs *ProcFS) parseDiskStatsLine(line string) (*BlkdevIOInfo, error) {
	var err error
	ret := &BlkdevIOInfo{}
	fields, err := procfs.SplitFields(line, 14)
	if err != nil {
		return ret, err
	}
	if ret.MajorNumber, err = procfs.ParseUint32(fields[0]); err != nil {
		return ret, err
	}
	if ret.MinorNumber, err = procfs.ParseUint32(fields[1]); err != nil {
		return ret, err
	}
	ret.DeviceName = fields[2]
	if ret.ReadsIOs, err = procfs.ParseUint64(fields[3]); err != nil {
		return ret, err
	}
	if ret.ReadsMerged, err = procfs.ParseUint64(fields[4]); err != nil {
		return ret, err
	}
	if ret.ReadsBytes, err = procfs.ParseMultUint64(fields[5], SectorSize); err != nil {
		return ret, err
	}
	if ret.ReadTimeMS, err = procfs.ParseUint64(fields[6]); err != nil {
		return ret, err
	}
	if ret.WritesIOs, err = procfs.ParseUint64(fields[7]); err != nil {
		return ret, err
	}
	if ret.WritesMerged, err = procfs.ParseUint64(fields[8]); err != nil {
		return ret, err
	}
	if ret.WritesBytes, err = procfs.ParseMultUint64(fields[9], SectorSize); err != nil {
		return ret, err
	}
	if ret.WriteTimeMS, err = procfs.ParseUint64(fields[10]); err != nil {
		return ret, err
	}
	if ret.InFlight, err = procfs.ParseUint64(fields[11]); err != nil {
		return ret, err
	}
	if ret.IOTimeMS, err = procfs.ParseUint64(fields[12]); err != nil {
		return ret, err
	}
	if ret.WeightedIOTimeMS, err = procfs.ParseUint64(fields[13]); err != nil {
		return ret, err
	}

	return ret, nil
}

// Devices parses "/proc/devices" into maps of character-devices and
// block-devices, indexed by device-number. See man(5) proc.
func (procfs *ProcFS) Devices() (map[int]string, map[int]string, error) {
	blkdev := make(map[int]string)
	chrdev := make(map[int]string)
	lines, err := procfs.ReadFileLines("devices")
	if err != nil {
		return blkdev, chrdev, err
	}
	inblk := false
	inchr := false
	for _, line := range lines {
		tline := strings.TrimSpace(line)
		if len(tline) == 0 {
			continue
		}
		if strings.HasPrefix(tline, "Character devices") {
			inchr = true

			continue
		}
		if strings.HasPrefix(tline, "Block devices") {
			inblk = true

			continue
		}
		dev, name, err := procfs.parseDevicesLine(tline)
		if err != nil {
			return blkdev, chrdev, err
		}
		if inchr {
			chrdev[dev] = name
		} else if inblk {
			blkdev[dev] = name
		}
	}

	return blkdev, chrdev, err
}

func (procfs *ProcFS) parseDevicesLine(line string) (int, string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, "", fmt.Errorf("bad line in devices: '%s'", line)
	}
	val, err := procfs.ParseInt(fields[0])
	if err != nil {
		return val, "", err
	}

	return val, fields[1], nil
}
