// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PseudoFS struct {
	Prefix string
}

func (pfs *PseudoFS) ParseBool(s string) (bool, error) {
	val, err := pfs.ParseInt(s)
	if err != nil {
		return false, err
	}

	return (val != 0), nil
}

func (pfs *PseudoFS) ParseInt(s string) (int, error) {
	val, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err == nil && val > math.MaxInt32 {
		err = fmt.Errorf("illegal int value: %s", s)
	}

	return int(val), err
}

func (pfs *PseudoFS) ParseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(strings.TrimSpace(s), 10, 32)
	if err == nil && val > math.MaxUint32 {
		err = fmt.Errorf("illegal uint32 value: %s", s)
	}

	return uint32(val), err
}

func (pfs *PseudoFS) ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 10, 64)
}

func (pfs *PseudoFS) ParseMultUint64(s string, mul uint64) (uint64, error) {
	ret, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return ret, err
	}
	return ret * mul, nil
}

func (pfs *PseudoFS) ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func (pfs *PseudoFS) ReadFile(subs ...string) (string, error) {
	return pfs.readFile(subs)
}

func (pfs *PseudoFS) ReadFileTrim(subs ...string) (string, error) {
	dat, err := pfs.readFile(subs)
	if err == nil {
		dat = strings.TrimSpace(dat)
	}
	return dat, err
}

func (pfs *PseudoFS) ReadFileAsBool(subs ...string) (bool, error) {
	dat, err := pfs.ReadFile(subs...)
	if err != nil {
		return false, err
	}
	return pfs.ParseBool(dat)
}

func (pfs *PseudoFS) ReadFileAsInt(subs ...string) (int, error) {
	dat, err := pfs.ReadFile(subs...)
	if err != nil {
		return 0, err
	}
	return pfs.ParseInt(dat)
}

func (pfs *PseudoFS) ReadFileAsUInt32(subs ...string) (uint32, error) {
	dat, err := pfs.ReadFile(subs...)
	if err != nil {
		return 0, err
	}
	return pfs.ParseUint32(dat)
}

func (pfs *PseudoFS) ReadFileAsUInt64(subs ...string) (uint64, error) {
	dat, err := pfs.ReadFile(subs...)
	if err != nil {
		return 0, err
	}
	return pfs.ParseUint64(dat)
}

func (pfs *PseudoFS) ReadFileLines(subs ...string) ([]string, error) {
	dat, err := pfs.readFile(subs)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(dat, "\n"), nil
}

func (pfs *PseudoFS) ReadFileFields(subs ...string) ([]string, error) {
	dat, err := pfs.readFile(subs)
	if err != nil {
		return []string{}, err
	}

	return strings.Fields(dat), nil
}

func (pfs *PseudoFS) readFile(subs []string) (string, error) {
	return readTextFile(pfs.resolvePath(subs))
}

func (pfs *PseudoFS) resolvePath(subs []string) string {
	return filepath.Join(pfs.Prefix, strings.Join(subs, "/"))
}

func (pfs *PseudoFS) ReadDir(subs ...string) ([]string, error) {
	return pfs.readDir(subs)
}

func (pfs *PseudoFS) readDir(subs []string) ([]string, error) {
	return pfs.listDir(pfs.resolvePath(subs))
}

func (pfs *PseudoFS) listDir(dirpath string) ([]string, error) {
	ret := []string{}
	entries, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return ret, err
	}
	for _, ent := range entries {
		ret = append(ret, filepath.Join(dirpath, ent.Name()))
	}

	return ret, nil
}

func (pfs *PseudoFS) SplitFields(s string, min int) ([]string, error) {
	fields := strings.Fields(s)
	if len(fields) < min {
		return fields, fmt.Errorf("only %d fields", len(fields))
	}
	return fields, nil
}

func (pfs *PseudoFS) IsDir(subs ...string) (bool, error) {
	fi, err := os.Stat(pfs.resolvePath(subs))
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// readWholeFile uses ioutil.ReadAll to read contents of entire file, but
// without relying on os.Stat for file's size, as many files in /proc and
// /sys report incorrect file sizes (either 0 or 4096). Limits the number
// read-bytes to 1-mega.
func readWholeFile(pathname string) ([]byte, error) {
	file, err := os.Open(pathname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(io.LimitReader(file, Mega))
}

func readTextFile(pathname string) (string, error) {
	dat, err := readWholeFile(pathname)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}
