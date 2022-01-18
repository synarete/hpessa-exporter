// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type SsaEntry struct {
	Title  string
	Values map[string]string
}

type SsaLogicalDrive struct {
	SsaEntry
}

type SsaPhysicalDrive struct {
	SsaEntry
}

type SsaArray struct {
	SsaEntry
	LogicalDrive  []*SsaLogicalDrive
	PhysicalDrive []*SsaPhysicalDrive
}

type SsaSlot struct {
	SsaEntry
	Array []*SsaArray
}

type SsaConfigInfo struct {
	Slots []*SsaSlot
}

type SsaPhysicalDriveInfo struct {
	ID         string `json:"id"`
	Box        string `json:"box"`
	Bay        string `json:"bay"`
	Size       string `json:"size"`
	SizeBytes  uint64 `json:"sizebytes"`
	Status     string `json:"status"`
	Serial     string `json:"serial"`
	TempCurr   int64  `json:"tempcurr"`
	TempMaxi   int64  `json:"tempmaxi"`
	UniqueID   string `json:"uniqueid"`
	PowerHours int64  `json:"powerhours"` // nolint:misspell
}

type SsaLogicalDriveInfo struct {
	ID             string `json:"id"`
	DiskName       string `json:"diskname"`
	Size           string `json:"size"`
	SizeBytes      uint64 `json:"sizebytes"`
	Status         string `json:"status"`
	UniqueID       string `json:"uniqueid"`
	ArrayName      string `json:"arrayname"`
	PhysicalDrives []SsaPhysicalDriveInfo
}

type SsaMap struct {
	DevMap map[string]SsaLogicalDriveInfo
}

func LocateSsa() (string, error) {
	knowns := []string{
		"/usr/sbin/ssacli",
		"/opt/smartstorageadmin/ssacli/bin/ssacli",
		"/opt/hp/ssacli/bld/ssacli",
	}
	for _, loc := range knowns {
		fi, err := os.Stat(loc)
		if err != nil {
			continue
		}
		mode := fi.Mode()
		if !mode.IsRegular() {
			continue
		}
		if (mode & 0111) > 0 {
			return loc, nil
		}
	}
	return "", errors.New("failed to locate ssacli")
}

func RunSsaVersion() (string, error) {
	dat, err := executeSsaCommand("version")
	if err != nil {
		return "", err
	}
	return ParseSsaVersion(dat)
}

func ParseSsaVersion(dat string) (string, error) {
	for _, ln := range strings.Split(dat, "\n") {
		sln := strings.Split(strings.TrimSpace(ln), ": ")
		if len(sln) < 2 {
			continue
		}
		if strings.HasPrefix(sln[0], "SSACLI Version") {
			return sln[1], nil
		}
	}
	return "", errors.New("failed to parse ssacli version")
}

func RunSsaShowConfig() (*SsaConfigInfo, error) {
	out, err := executeSsaCommand("ctrl", "all", "show", "config", "detail")
	if err != nil {
		return nil, err
	}
	return ParseSsaShowConfig(out)
}

func executeSsaCommand(args ...string) (string, error) {
	loc, err := LocateSsa()
	if err != nil {
		return "", err
	}
	return executeCommand(loc, args...)
}

func ParseSsaShowConfig(dat string) (*SsaConfigInfo, error) {
	cfg := &SsaConfigInfo{}
	var slot *SsaSlot
	var arr *SsaArray
	var ld *SsaLogicalDrive
	var pd *SsaPhysicalDrive
	haskv := false
	subsec := false
	for _, line := range strings.Split(dat, "\n") {
		ln, key, value, _, newsec := parseSsaLine(line)
		haskv = len(key) > 0 && len(value) > 0
		if newsec {
			ld = nil
			pd = nil
			subsec = true
		}

		if strings.Index(ln, "Slot ") > 0 {
			appendSlotInfo(cfg, slot)
			slot = newSsaSlot(ln)
			subsec = false
			arr = nil
			ld = nil
			pd = nil
			continue
		}

		if slot == nil || hasIgnoredSubSections(line) {
			arr = nil
			ld = nil
			pd = nil
			continue
		}

		switch {
		case hasSubSections(line, "Array: "):
			arr = newSsaArray(ln)
			slot.Array = append(slot.Array, arr)
		case arr != nil && hasSubSections(line, "Logical Drive: "):
			pd = nil
			ld = newSsaLogicalDrive(ln)
			arr.LogicalDrive = append(arr.LogicalDrive, ld)
		case arr != nil && hasSubSections(line, "physicaldrive "):
			ld = nil
			pd = newSsaPhysicalDrive(ln)
			arr.PhysicalDrive = append(arr.PhysicalDrive, pd)
		case pd != nil && haskv:
			pd.Values[key] = value
		case ld != nil && haskv:
			ld.Values[key] = value
		case arr != nil && haskv:
			arr.Values[key] = value
		case slot != nil && haskv && !subsec:
			slot.Values[key] = value
		default:
		}
	}
	appendSlotInfo(cfg, slot)
	return cfg, nil
}

func hasIgnoredSubSections(line string) bool {
	return hasSubSections(line, "Internal Drive Cage", "Physical Drives", "Port Name")
}

func hasSubSections(line string, secs ...string) bool {
	if countLeading(line, ' ') == 0 {
		return false
	}
	for _, sec := range secs {
		sline := strings.TrimSpace(line)
		if strings.Index(sline, sec) == 0 {
			subs := strings.Split(sline, " ")
			if len(subs) > 2 && subs[0] == "physicaldrive" {
				return false
			}
			return true
		}
	}
	return false
}

func newSsaSlot(title string) *SsaSlot {
	slot := &SsaSlot{}
	slot.Title = title
	slot.Values = map[string]string{}
	slot.Array = []*SsaArray{}
	return slot
}

func newSsaArray(title string) *SsaArray {
	array := &SsaArray{}
	array.Title = title
	array.Values = map[string]string{}
	array.LogicalDrive = []*SsaLogicalDrive{}
	array.PhysicalDrive = []*SsaPhysicalDrive{}
	return array
}

func newSsaLogicalDrive(title string) *SsaLogicalDrive {
	ld := &SsaLogicalDrive{}
	ld.Title = title
	ld.Values = map[string]string{}
	return ld
}

func newSsaPhysicalDrive(title string) *SsaPhysicalDrive {
	pd := &SsaPhysicalDrive{}
	pd.Title = title
	pd.Values = map[string]string{}
	return pd
}

func appendSlotInfo(cfg *SsaConfigInfo, slotInfo *SsaSlot) {
	if slotInfo != nil {
		cfg.Slots = append(cfg.Slots, slotInfo)
	}
}

func parseSsaLine(line string) (string, string, string, int, bool) {
	ln := strings.TrimSpace(line)
	if len(ln) == 0 {
		return "", "", "", 0, true
	}
	ind := countLeading(line, ' ')
	key, val := keyValueOf(ln)
	if len(key) > 0 && len(val) > 0 {
		return ln, key, val, ind, false
	}
	return ln, "", "", ind, false
}

func countLeading(str string, item rune) int {
	cnt := 0
	for _, val := range str {
		if val == item {
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func ParseConfigToLogical(config *SsaConfigInfo) (*SsaMap, error) {
	ret := NewSsaMap()
	for _, slot := range config.Slots {
		for _, arr := range slot.Array {
			if len(arr.LogicalDrive) != 1 {
				continue
			}
			ld := arr.LogicalDrive[0]
			name := path.Base(ld.Values["Disk Name"])
			if name == "" {
				continue
			}
			ldi := SsaLogicalDriveInfo{}
			ldi.ID = valueOf(ld.Title)
			ldi.DiskName = ld.Values["Disk Name"]
			ldi.Size = ld.Values["Size"]
			ldi.SizeBytes = parseSizeBytes(ldi.Size)
			ldi.Status = ld.Values["Status"]
			ldi.UniqueID = ld.Values["Unique Identifier"]
			ldi.ArrayName = valueOf(arr.Title)
			for _, pd := range arr.PhysicalDrive {
				pdi := SsaPhysicalDriveInfo{}
				pdi.ID = pd.Title
				pdi.Box = pd.Values["Box"]
				pdi.Bay = pd.Values["Bay"]
				pdi.Size = pd.Values["Size"]
				pdi.SizeBytes = parseSizeBytes(pdi.Size)
				pdi.Status = pd.Values["Status"]
				pdi.TempCurr = parseTemp(valueByPrefix(pd.Values, "Current Temperature"))
				pdi.TempMaxi = parseTemp(valueByPrefix(pd.Values, "Maximum Temperature"))
				pdi.UniqueID = pd.Values["Drive Unique ID"]
				pdi.PowerHours = parseInt64(pd.Values["Power On Hours"])

				ldi.PhysicalDrives = append(ldi.PhysicalDrives, pdi)
			}
			ret.DevMap[name] = ldi
		}
	}
	return ret, nil
}

func valueByPrefix(kv map[string]string, prefix string) string {
	for key, val := range kv {
		if strings.HasPrefix(key, prefix) {
			return val
		}
	}
	return ""
}

func valueOf(s string) string {
	_, v := keyValueOf(s)
	return v
}

func keyValueOf(s string) (string, string) {
	kv := strings.Split(s, ": ")
	if len(kv) != 2 {
		return strings.TrimSpace(s), ""
	}
	return strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
}

func parseSizeBytes(s string) uint64 {
	kv := strings.Split(s, " ")
	if len(kv) == 0 {
		return 0
	}
	val, err := strconv.ParseFloat(kv[0], 64)
	if err != nil {
		return 0
	}
	if len(kv) == 2 {
		switch strings.ToUpper(kv[1]) {
		case "KB":
			val *= float64(Kilo)
		case "GB":
			val *= float64(Giga)
		case "TB":
			val *= float64(Tera)
		case "PB":
			val *= float64(Peta)
		}
	}
	return uint64(val)
}

func parseTemp(s string) int64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(val)
}

func parseInt64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return val
}

func NewSsaMap() *SsaMap {
	return &SsaMap{
		DevMap: map[string]SsaLogicalDriveInfo{},
	}
}

func executeCommand(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}
	res := strings.TrimSpace(string(out))
	return res, nil
}
