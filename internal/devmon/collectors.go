// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	collectorsNamespace = "hpessa"
)

func (dex *deviceExporter) register() error {
	return dex.registerCollectors(dex.listCollectors())
}

func (dex *deviceExporter) registerCollectors(cols []prometheus.Collector) error {
	for _, c := range cols {
		if err := dex.reg.Register(c); err != nil {
			dex.log.Error(err, "failed to register collector")
			return err
		}
	}
	return nil
}

func (dex *deviceExporter) listCollectors() []prometheus.Collector {
	cols := []prometheus.Collector{}
	if dex.any {
		prc := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
		goc := collectors.NewGoCollector()
		cols = append(cols, prc, goc)
	}
	cols = append(cols, newExporterVersionCollector())
	cols = append(cols, dex.newSsaVersionCollector())
	cols = append(cols, dex.newBlkdevCollector())
	cols = append(cols, dex.newSsaLogicalDrivesCollector())
	cols = append(cols, dex.newSsaPhysicalDrivesCollector())
	return cols
}

func collectorName(subsystem, name string) string {
	return prometheus.BuildFQName(collectorsNamespace, subsystem, name)
}

func newExporterVersionCollector() prometheus.Collector {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: collectorName("exporter", "info"),
		Help: "Version of the exporter.",
		ConstLabels: map[string]string{
			"version": Version(),
		},
	})
	gauge.Set(1)
	return gauge
}

type deCollector struct {
	// nolint:structcheck
	dex *deviceExporter
	dsc []*prometheus.Desc
}

func (col *deCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range col.dsc {
		ch <- d
	}
}

type ssaVersionCollector struct {
	deCollector
}

func (col *ssaVersionCollector) Collect(ch chan<- prometheus.Metric) {
	vers, err := RunSsaVersion()
	if err != nil {
		vers = "N/A"
	}
	ch <- prometheus.MustNewConstMetric(col.dsc[0],
		prometheus.GaugeValue, 1, vers)
}

func (dex *deviceExporter) newSsaVersionCollector() prometheus.Collector {
	col := &ssaVersionCollector{}
	col.dex = dex
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("ssacli", "info"),
			"Version of local SSA cli utility.",
			[]string{"version"}, nil),
	}
	return col
}

type blkdevCollector struct {
	deCollector
}

func (col *blkdevCollector) Collect(ch chan<- prometheus.Metric) {
	bdis, _ := col.dex.sdp.probeBlockDevices()
	for _, bdi := range bdis {
		ch <- prometheus.MustNewConstMetric(col.dsc[0],
			prometheus.GaugeValue,
			float64(bdi.Size),
			bdi.Name,
			strconv.Itoa(int(bdi.Major)),
			strconv.Itoa(int(bdi.Minor)),
			bdi.Vendor,
			bdi.Model)
	}
}

func (dex *deviceExporter) newBlkdevCollector() prometheus.Collector {
	col := &blkdevCollector{}
	col.dex = dex
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("blkdev", "size_bytes"),
			"Block device size in bytes.",
			[]string{"name", "major", "minor", "vendor", "model"}, nil),
	}
	return col
}

type ssaLogicalDrivesCollector struct {
	deCollector
}

func (col *ssaLogicalDrivesCollector) Collect(ch chan<- prometheus.Metric) {
	sdis, _ := col.dex.sdp.probeDevices()
	for _, sdi := range sdis {
		ldi := sdi.SsaLogicalDrive
		if ldi == nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(col.dsc[0],
			prometheus.GaugeValue,
			statusToValue(ldi.Status),
			ldi.ArrayName, ldi.DiskName, ldi.Status)
	}
}

func (dex *deviceExporter) newSsaLogicalDrivesCollector() prometheus.Collector {
	col := &ssaLogicalDrivesCollector{}
	col.dex = dex
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("ssa_logical_device", "status"),
			"Status of logical device",
			[]string{"arrayname", "diskname", "status"}, nil),
	}
	return col
}

type ssaPhysicalDrivesCollector struct {
	deCollector
}

func (col *ssaPhysicalDrivesCollector) Collect(ch chan<- prometheus.Metric) {
	sdis, _ := col.dex.sdp.probeDevices()
	for _, sdi := range sdis {
		ldi := sdi.SsaLogicalDrive
		if ldi == nil {
			continue
		}
		for _, pdi := range ldi.PhysicalDrives {
			ch <- prometheus.MustNewConstMetric(col.dsc[0],
				prometheus.GaugeValue,
				statusToValue(pdi.Status),
				sdi.Name, pdi.ID, pdi.Box, pdi.Bay, pdi.UniqueID)

			ch <- prometheus.MustNewConstMetric(col.dsc[1],
				prometheus.GaugeValue,
				float64(pdi.TempCurr),
				sdi.Name, pdi.ID, pdi.Box, pdi.Bay, pdi.UniqueID)

			ch <- prometheus.MustNewConstMetric(col.dsc[2],
				prometheus.GaugeValue,
				float64(pdi.TempMaxi),
				sdi.Name, pdi.ID, pdi.Box, pdi.Bay, pdi.UniqueID)
		}
	}
}

func (dex *deviceExporter) newSsaPhysicalDrivesCollector() prometheus.Collector {
	col := &ssaPhysicalDrivesCollector{}
	col.dex = dex
	col.dsc = []*prometheus.Desc{
		prometheus.NewDesc(
			collectorName("ssa_physical_device", "status"),
			"Status of physical device",
			[]string{"dev", "id", "box", "bay", "uniqueid"}, nil),
		prometheus.NewDesc(
			collectorName("ssa_physical_device", "temp_curr"),
			"Current temperature of physical device",
			[]string{"dev", "id", "box", "bay", "uniqueid"}, nil),
		prometheus.NewDesc(
			collectorName("ssa_physical_device", "temp_maxi"),
			"Maximal temperature of physical device",
			[]string{"dev", "id", "box", "bay", "uniqueid"}, nil),
	}
	return col
}

func statusToValue(status string) float64 {
	return float64(statusToInt(status))
}

func statusToInt(status string) int {
	var ret int
	if strings.ToUpper(status) != "OK" {
		ret = 1
	} else {
		ret = 0
	}
	return ret
}
