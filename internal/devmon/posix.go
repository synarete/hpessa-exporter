// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"bytes"

	"golang.org/x/sys/unix"
)

func Uname() (*SysUname, error) {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return &SysUname{}, err
	}

	return &SysUname{
		Sysname:  unameBytesToString(uname.Sysname),
		Nodename: unameBytesToString(uname.Nodename),
		Release:  unameBytesToString(uname.Release),
		Version:  unameBytesToString(uname.Version),
		Machine:  unameBytesToString(uname.Machine),
	}, nil
}

func unameBytesToString(b [65]byte) string {
	n := bytes.IndexByte(b[:], 0)
	if n < 0 {
		n = len(b)
	}

	return string(b[:n])
}
