// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"fmt"
	"net"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	DefaultMetricsPort = int(8080)
)

type deviceExporter struct {
	log  logr.Logger
	sdp  *storageDevicesProbe
	reg  *prometheus.Registry
	mux  *http.ServeMux
	port int
	any  bool
}

func newDeviceExporter(log logr.Logger, port int) *deviceExporter {
	return &deviceExporter{
		log:  log,
		sdp:  newStorageDevicesProbe(log),
		reg:  prometheus.NewRegistry(),
		mux:  http.NewServeMux(),
		port: port,
		any:  false,
	}
}

func (dex *deviceExporter) init() error {
	dex.log.Info("init devices exporter")
	if err := dex.sdp.init(); err != nil {
		return err
	}
	dex.log.Info("register collectors")
	if err := dex.register(); err != nil {
		return err
	}
	return nil
}

func (dex *deviceExporter) serve() error {
	addr := fmt.Sprintf(":%d", dex.port)
	dex.log.Info("serve metrics", "addr", addr)

	handler := promhttp.HandlerFor(dex.reg, promhttp.HandlerOpts{})
	dex.mux.Handle("/metrics", handler)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		dex.log.Error(err, "failed to listen", "addr", addr)
		return err
	}
	defer listener.Close()

	if err := http.Serve(listener, dex.mux); err != nil {
		dex.log.Error(err, "HTTP server failure", "addr", addr)
		return err
	}
	return nil
}

func RunDevicesExporter(port int) error {
	log := zap.New(zap.UseFlagOptions(&zap.Options{}))
	dex := newDeviceExporter(log, port)
	if err := dex.init(); err != nil {
		return err
	}
	if err := dex.serve(); err != nil {
		return err
	}
	return nil
}

func ProbePrintDevices() error {
	log := zap.New(zap.UseFlagOptions(&zap.Options{}))
	sdp := newStorageDevicesProbe(log)

	if err := sdp.init(); err != nil {
		return err
	}
	sdi, err := sdp.probeDevices()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", sdi)
	return nil
}
