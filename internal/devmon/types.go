// SPDX-License-Identifier: Apache-2.0
package devmon

type SysUname struct {
	Sysname  string `json:"sysname"`
	Nodename string `json:"nodename"`
	Release  string `json:"release"`
	Version  string `json:"version"`
	Machine  string `json:"machine"`
}

type NodeInfo struct {
	SysUname
	HostIP []string `json:"hostip"`
}

type SysLoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// BlkdevInfo contains collection of raw information under /sys/block/<disk>/...
// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-block
type BlkdevInfo struct {
	Major    uint32 `json:"major"`
	Minor    uint32 `json:"minor"`
	Name     string `json:"name"`
	Size     uint64 `json:"size"`
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Readonly bool   `json:"readonly"`
}

type BlkdevID struct {
	MajorNumber uint32 `json:"major"`
	MinorNumber uint32 `json:"minor"`
	DeviceName  string `json:"devname"`
}

// BlkdevIOStat represents procfs diskstats info
// https://www.kernel.org/doc/Documentation/iostats.txt,
// https://www.kernel.org/doc/Documentation/block/stat.txt,
// https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
type BlkdevIOStat struct {
	ReadsIOs         uint64 `json:"readios"`
	ReadsMerged      uint64 `json:"readsmerged"`
	ReadsBytes       uint64 `json:"readsbytes"`
	ReadTimeMS       uint64 `json:"readtimems"`
	WritesIOs        uint64 `json:"writeios"`
	WritesMerged     uint64 `json:"writesmerged"`
	WritesBytes      uint64 `json:"writesbytes"`
	WriteTimeMS      uint64 `json:"writetimems"`
	InFlight         uint64 `json:"inflight"`
	IOTimeMS         uint64 `json:"iotimems"`
	WeightedIOTimeMS uint64 `json:"weightediotimems"`
}

type BlkdevIOInfo struct {
	BlkdevID
	BlkdevIOStat
}

// BlkdevQueueStats represents sysfs block-queue info
// https://www.kernel.org/doc/Documentation/block/queue-sysfs.txt
// https://www.kernel.org/doc/html/latest/block/queue-sysfs.html
type BlkdevQueueStats struct {
	AddRandom            bool   `json:"addrandom"`
	DAX                  bool   `json:"dax"`
	DiscardGranularity   uint64 `json:"discardgranularity"`
	DiscardMaxHWBytes    uint64 `json:"discardmaxhwbytes"`
	DiscardMaxBytes      uint64 `json:"discardmaxbytes"`
	FUA                  bool   `json:"fua"`
	HWSectorSize         uint32 `json:"hwsectorsize"`
	IOPoll               bool   `json:"iopoll"`
	IOPollDelay          int    `json:"iopolldelay"`
	IOTimeout            uint64 `json:"iotimeout"`
	IOStats              bool   `json:"iostats"`
	LogicalBlockSize     uint64 `json:"logicalblocksize"`
	MaxHWSectorsKB       uint64 `json:"maxhwsectorskb"`
	MaxIntegritySegments uint64 `json:"maxintegritysegments"`
	MaxSectorsKB         uint64 `json:"maxsectorskb"`
	MaxSegments          uint64 `json:"maxsegments"`
	MaxSegmentSize       uint64 `json:"maxsegmentsize"`
	MinimumIOSize        uint64 `json:"minimumiosize"`
	NoMerges             uint32 `json:"nomerges"`
	NRRequests           uint64 `json:"nrrequests"`
	OptimalIOSize        uint64 `json:"optimaliosize"`
	PhysicalBlockSize    uint64 `json:"physicalblocksize"`
	ReadAHeadKB          uint64 `json:"readaheadkb"`
	Rotational           bool   `json:"rotational"`
	RQAffinity           uint32 `json:"rqaffinity"`
	Scheduler            string `json:"scheduler"`
	WriteCache           string `json:"writecache"`
	WriteSameMaxBytes    uint64 `json:"writesamemaxbytes"`
	WBTLatUSec           int    `json:"wbtlatusec"`
	Zoned                string `json:"zoned"`
	ZoneWriteGranularity int    `json:"zonewritegranularity"`
}
