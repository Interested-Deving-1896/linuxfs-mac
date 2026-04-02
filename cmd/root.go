// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Global flags shared across subcommands.
var (
	flagVMMemMiB  uint
	flagDebug     bool
	flagDataDir   string
	flagBackend   string
	flagListenIP  string
)

func init() {
	flag.UintVar(&flagVMMemMiB, "vm-mem", 512,
		"RAM allocated to the Alpine VM in MiB (use >=2048 for LUKS)")
	flag.BoolVar(&flagDebug, "debug", false,
		"Enable verbose VM and QEMU output")
	flag.StringVar(&flagDataDir, "data-dir", "",
		"Directory for VM image and state (default: OS cache dir)")
	flag.StringVar(&flagBackend, "backend", "",
		"File share backend: smb, afp, nfs, ftp (auto-detected by OS if empty)")
	flag.StringVar(&flagListenIP, "listen-ip", "127.0.0.1",
		"IP address the share server listens on")
}

// Execute parses args and dispatches to the appropriate subcommand.
func Execute() {
	flag.Usage = usage
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	sub := os.Args[1]
	// Remaining args after the subcommand name.
	rest := os.Args[2:]

	switch sub {
	case "mount":
		runMount(rest)
	case "unmount":
		runUnmount(rest)
	case "list":
		runList(rest)
	case "shell":
		runShell(rest)
	case "version":
		runVersion()
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %q\n\n", sub)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `linuxfs-mac — access Linux-native filesystems on macOS and Windows.

Usage:
  linuxfs-mac <subcommand> [flags]

Subcommands:
  mount     Mount a Linux filesystem and expose it as a network share
  unmount   Unmount a previously mounted filesystem
  list      List currently mounted filesystems
  shell     Open a shell inside the Alpine VM for a device
  version   Print version information

Global flags:
`)
	flag.PrintDefaults()
}
