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
	showVersion    bool
	probePrintOnce bool

	rootCmd = &cobra.Command{
		Use:   "hpessa-exporter",
		Short: "Storage devices monitor",
		Run:   func(cmd *cobra.Command, args []string) { start() },
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&showVersion,
		"version", "v", false, "show version and exit")
	rootCmd.Flags().BoolVarP(&probePrintOnce,
		"print", "p", false, "probe-print devices and exit")
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
	if probePrintOnce {
		if err := devmon.ProbePrintDevices(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := devmon.RunDevicesExporter(); err != nil {
		os.Exit(1)
	}
}
