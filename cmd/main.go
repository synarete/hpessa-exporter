// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/spf13/cobra"
)

var (
	showVersion bool
	showDevices bool
	metricsPort int

	rootCmd = &cobra.Command{
		Use:   "hpessa-exporter",
		Short: "Storage devices monitor",
		Run:   func(cmd *cobra.Command, args []string) { start() },
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&showVersion,
		"version", "v", false, "show version and exit")
	rootCmd.Flags().BoolVarP(&showDevices,
		"show", "s", false, "probe-print devices and exit")
	rootCmd.Flags().IntVarP(&metricsPort,
		"port", "p", devmon.DefaultMetricsPort, "metrics port")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func start() {
	if showVersion {
		fmt.Printf("%s %s %s %s\n", devmon.Version(),
			runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	if showDevices {
		if err := devmon.ProbePrintDevices(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := devmon.RunDevicesExporter(metricsPort); err != nil {
		os.Exit(1)
	}
}
