// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// storageDeviceInfo represents a union of information collected on a block
// storage device: both from the Linux platform and optionally, via external
// vendor-specific tools.
type storageDeviceInfo struct {
	BlkdevInfo
	SsaLogicalDrive *SsaLogicalDriveInfo
}

// storageDevicesProbe is an auxiliary object to collect storage-devices info
// from system and externa-tools.
type storageDevicesProbe struct {
	log    logr.Logger
	ident  *Ident
	procfs *ProcFS
	sysfs  *SysFS
	clnt   *client
	hasSSA bool
}

func newStorageDevicesProbe(log logr.Logger) *storageDevicesProbe {
	return &storageDevicesProbe{
		log:    log,
		ident:  SelfIdent(),
		procfs: NewProcFS(),
		sysfs:  NewSysFS(),
		hasSSA: true,
	}
}

func (sdp *storageDevicesProbe) init() error {
	if err := sdp.initClient(); err != nil {
		return err
	}
	if err := sdp.initSelf(); err != nil {
		return err
	}
	if err := sdp.initXtool(); err != nil {
		return err
	}
	return nil
}

func (sdp *storageDevicesProbe) initXtool() error {
	loc, err := LocateSsa()
	if err != nil {
		sdp.hasSSA = false
		sdp.log.Info("unable to find ssacli tool")
		return nil // OK
	}
	sdp.log.Info("found ssacli tool at: " + loc)
	vers, err := RunSsaVersion()
	if err != nil {
		sdp.hasSSA = false
		sdp.log.Error(err, "failed to run ssacli tool")
		return err
	}
	sdp.log.Info("ssacli", "version", vers)
	return nil
}

func (sdp *storageDevicesProbe) initClient() error {
	kclnt, err := newClient()
	if err != nil {
		sdp.log.Error(err, "failed to create clientset")
		return err
	}
	sdp.clnt = kclnt
	return nil
}

func (sdp *storageDevicesProbe) initSelf() error {
	ident := SelfIdent()
	if ident.HostIP == "" {
		return fmt.Errorf("unable to resolve hostip: %+v", ident)
	}
	pod, err := sdp.discoverSelfPod()
	if err != nil {
		return err
	}
	sdp.log.Info("self pod", "name", pod.GetName(),
		"podip", pod.Status.PodIP, "ident", ident)
	return nil
}

func (sdp *storageDevicesProbe) probeDevices() ([]storageDeviceInfo, error) {
	bdi, err := sdp.probeBlockDevices()
	if err != nil {
		return []storageDeviceInfo{}, err
	}
	ssm, err := sdp.probeSsaLogicalDevices()
	if err != nil {
		return []storageDeviceInfo{}, err
	}
	return newStorageDeviceInfo(bdi, ssm), nil
}

func newStorageDeviceInfo(bdis []BlkdevInfo, ssm *SsaMap) []storageDeviceInfo {
	ret := []storageDeviceInfo{}
	for _, bdi := range bdis {
		sdi := storageDeviceInfo{
			BlkdevInfo: bdi,
		}
		ssaent, ok := ssm.DevMap[bdi.Name]
		if ok {
			sdi.SsaLogicalDrive = &ssaent
		}
		ret = append(ret, sdi)
	}
	return ret
}

func (sdp *storageDevicesProbe) probeBlockDevices() ([]BlkdevInfo, error) {
	bdi, err := DiscoverBlkdevInfo(sdp.procfs, sdp.sysfs)
	if err != nil {
		sdp.log.Error(err, "failed to discover block devices")
		return []BlkdevInfo{}, err
	}
	return bdi, nil
}

func (sdp *storageDevicesProbe) probeBlockDevicesIO() ([]BlkdevIOInfo, error) {
	ret, err := DescoveBlockDevicesIO(sdp.procfs, sdp.sysfs)
	if err != nil {
		sdp.log.Error(err, "failed to discover block devices IO stats")
		return []BlkdevIOInfo{}, err
	}
	return ret, nil
}

func (sdp *storageDevicesProbe) probeSsaLogicalDevices() (*SsaMap, error) {
	ssm := NewSsaMap()
	if _, err := LocateSsa(); err != nil {
		return ssm, nil // OK -- run without ssacli
	}
	if _, err := RunSsaVersion(); err != nil {
		return ssm, err
	}
	cfg, err := RunSsaShowConfig()
	if err != nil {
		sdp.log.Error(err, "failed to run ssacli show config")
		return nil, err
	}

	ldm, err := ParseConfigToLogical(cfg)
	if err != nil {
		sdp.log.Error(err, "failed to parse ssacli show config output")
		return nil, err
	}

	ssm, err = sdp.filterDevices(ldm)
	if err != nil {
		sdp.log.Error(err, "failed to filter logical devices")
		return nil, err
	}
	return ssm, nil
}

func (sdp *storageDevicesProbe) filterDevices(ldm *SsaMap) (*SsaMap, error) {
	ssm := NewSsaMap()
	bdm, err := DiscoverBlkdevStats(sdp.procfs, sdp.sysfs)
	if err != nil {
		sdp.log.Error(err, "failed to discover block devices")
		return ssm, err
	}
	for dev, val := range ldm.DevMap {
		if bdm.IOStats[dev] != nil {
			ssm.DevMap[dev] = val
		}
	}
	return ssm, nil
}

func (sdp *storageDevicesProbe) discoverSelfPod() (*corev1.Pod, error) {
	if sdp.clnt == nil {
		return nil, errors.New("no kube client")
	}
	nname := types.NamespacedName{
		Namespace: sdp.ident.Namespace,
		Name:      sdp.ident.Name,
	}
	pod, err := GetRunningPod(context.TODO(), sdp.clnt, nname)
	if err != nil {
		sdp.log.Error(err, "failed to get self pod")
		return nil, err
	}
	return pod, nil
}
