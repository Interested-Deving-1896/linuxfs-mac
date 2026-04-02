// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linuxfs-mac",
	Short: "Access Linux-native filesystems (ext4, btrfs, XFS, LVM, LUKS) on macOS and Windows.",
	Long: `linuxfs-mac mounts Linux filesystems on macOS and Windows without reimplementing
any filesystem driver. It spins up a lightweight Alpine Linux VM, passes through
the target block device, mounts it natively inside the VM, then exposes it to
the host via SMB (Windows), AFP or NFS (macOS), or FTP (fallback).

Supports: ext2/3/4, btrfs, XFS, ZFS, NTFS, exFAT, LVM, LUKS, BitLocker, RAID.`,
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Global flags shared across subcommands.
var (
	flagVMMemMiB   uint32
	flagDebug      bool
	flagDataDir    string
	flagBackend    string
	flagListenIP   string
)

func init() {
	rootCmd.PersistentFlags().Uint32Var(&flagVMMemMiB, "vm-mem", 512,
		"RAM allocated to the Alpine VM in MiB (use ≥2048 for LUKS)")
	rootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false,
		"Enable verbose VM and QEMU output")
	rootCmd.PersistentFlags().StringVar(&flagDataDir, "data-dir", "",
		"Directory for VM image and state (default: OS cache dir)")
	rootCmd.PersistentFlags().StringVar(&flagBackend, "backend", "",
		"File share backend: smb, afp, nfs, ftp (auto-detected by OS if empty)")
	rootCmd.PersistentFlags().StringVar(&flagListenIP, "listen-ip", "127.0.0.1",
		"IP address the share server listens on")

	rootCmd.AddCommand(mountCmd)
	rootCmd.AddCommand(unmountCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(versionCmd)
}
